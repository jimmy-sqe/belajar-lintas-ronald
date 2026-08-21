package app

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"

	// boilerplate:axis=auth option=jwt START
	jwtadapter "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/auth/jwt"
	pgpersist "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/persistence/postgres"
	rediscache "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/cache/redis"
	jwtissuer "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/token/jwt"
	redissession "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/session/redis"
	// boilerplate:axis=auth option=jwt END

	authdomain "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/auth"

	// boilerplate:axis=rpc option=grpc START
	grpcadapter "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/rpc/grpc"
	healthv1 "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/gen/health/v1"
	// boilerplate:axis=rpc option=grpc END

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/storage"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/config"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/inference"
	// boilerplate:axis=sample-app option=todo-app START
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/todo"
	// boilerplate:axis=sample-app option=todo-app END
	tehttp "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/http"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/http/handler"
)

// Container holds long-lived application dependencies.
type Container struct {
	Cfg *config.Config
	// boilerplate:axis=sample-app option=todo-app START
	TodoSvc todo.Service
	// boilerplate:axis=sample-app option=todo-app END
	HTTPHandlers   *tehttp.Handlers
	HTTPMiddleware *tehttp.MiddlewareStack
	Storage        storage.Storage
	Inference      inference.Engine
	Validator      echo.Validator
	ObsMiddleware  echo.MiddlewareFunc

	// boilerplate:axis=rpc option=grpc START
	GRPCServer *grpcadapter.Server
	// boilerplate:axis=rpc option=grpc END

	shutdowns []func(context.Context) error
}

// NewContainer constructs the application container.
func NewContainer(ctx context.Context, cfg *config.Config) (*Container, error) {
	c := &Container{Cfg: cfg, Validator: tehttp.NewValidator()}

	// === Persistence + Cache (sample-app) ===
	// boilerplate:axis=sample-app option=todo-app START
	todoRepo, todoRepoClose, err := newTodoRepo(ctx, cfg)
	if err != nil {
		return nil, err
	}
	c.shutdowns = append(c.shutdowns, todoRepoClose)

	todoCache, todoCacheClose, err := newTodoCache(ctx, cfg)
	if err != nil {
		return nil, err
	}
	c.shutdowns = append(c.shutdowns, todoCacheClose)
	// boilerplate:axis=sample-app option=todo-app END

	// === Auth service (postgres + redis pinned) ===
	// boilerplate:axis=auth option=jwt START
	authPgDB, err := pgpersist.Connect(ctx, cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("app: auth postgres: %w", err)
	}
	c.shutdowns = append(c.shutdowns, func(_ context.Context) error { return authPgDB.Close() })
	userRepo := pgpersist.NewUserRepository(authPgDB)

	authRDB, err := rediscache.Connect(ctx, cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("app: auth redis: %w", err)
	}
	c.shutdowns = append(c.shutdowns, func(_ context.Context) error { return authRDB.Close() })
	sessionStore := redissession.NewSessionStoreWithClient(authRDB)

	tokenIssuer := jwtissuer.NewIssuer(cfg.JWT.Secret)
	authenticator := authdomain.NewService(userRepo, sessionStore, tokenIssuer, authdomain.ServiceConfig{
		AccessTTL:  time.Duration(cfg.JWT.AccessTTLSec) * time.Second,
		RefreshTTL: time.Duration(cfg.JWT.RefreshTTLSec) * time.Second,
	})
	// boilerplate:axis=auth option=jwt END

	// === Service ===
	// boilerplate:axis=sample-app option=todo-app START
	c.TodoSvc = todo.NewService(todoRepo, todoCache)
	// boilerplate:axis=sample-app option=todo-app END

	// === Auth middleware ===
	mw := &tehttp.MiddlewareStack{
		Recover:   echomw.Recover(),
		RequestID: echomw.RequestID(),
	}
	// boilerplate:axis=auth option=jwt START
	jwtVerifier, err := jwtadapter.NewVerifier(cfg.JWT)
	if err != nil {
		return nil, fmt.Errorf("app: jwt: %w", err)
	}
	mw.Auth = jwtadapter.Middleware(jwtVerifier)
	// boilerplate:axis=auth option=jwt END

	c.HTTPMiddleware = mw

	// === Observability ===
	obs, err := newObservability(ctx, cfg)
	if err != nil {
		return nil, err
	}
	c.shutdowns = append(c.shutdowns, obs.shutdown)
	c.ObsMiddleware = obs.middleware

	// === Storage ===
	storageSvc, storageShutdown, err := newStorage(ctx, cfg)
	if err != nil {
		return nil, err
	}
	c.Storage = storageSvc
	if storageShutdown != nil {
		c.shutdowns = append(c.shutdowns, storageShutdown)
	}

	// === Inference ===
	inferenceEngine, err := newInference(cfg)
	if err != nil {
		return nil, err
	}
	c.Inference = inferenceEngine

	// boilerplate:axis=rpc option=grpc START
	grpcServer := grpcadapter.NewServer(cfg.GRPCHost, cfg.GRPCPort)
	healthv1.RegisterHealthServer(grpcServer.Registrar(), grpcadapter.NewHealthService())
	c.GRPCServer = grpcServer
	// boilerplate:axis=rpc option=grpc END

	// === Handlers ===
	c.HTTPHandlers = &tehttp.Handlers{
		// boilerplate:axis=auth option=jwt START
		Auth: handler.NewAuthHandler(authenticator),
		// boilerplate:axis=auth option=jwt END
		// boilerplate:axis=sample-app option=todo-app START
		Todo: handler.NewTodoHandler(c.TodoSvc),
		// boilerplate:axis=sample-app option=todo-app END
	}
	return c, nil
}

// Close releases all resources in reverse order of acquisition.
func (c *Container) Close() error {
	ctx := context.Background()
	for i := len(c.shutdowns) - 1; i >= 0; i-- {
		_ = c.shutdowns[i](ctx)
	}
	return nil
}

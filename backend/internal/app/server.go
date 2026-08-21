package app

import (
	"context"
	"fmt"
	// boilerplate:axis=rpc option=grpc START
	"log"
	// boilerplate:axis=rpc option=grpc END
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	tehttp "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/http"
)

// Server wraps an Echo HTTP server bound to the container.
type Server struct {
	echo      *echo.Echo
	container *Container
}

// NewServer builds a Server with routes and middleware registered.
func NewServer(c *Container) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Validator = c.Validator

	if len(c.Cfg.CORS.AllowedOrigins) > 0 {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     c.Cfg.CORS.AllowedOrigins,
			AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
			AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderXRequestedWith, "X-Utc-Offset"},
			AllowCredentials: true,
			MaxAge:           86400,
		}))
	}

	e.Use(c.ObsMiddleware)

	tehttp.RegisterRoutes(e, c.HTTPHandlers, c.HTTPMiddleware)

	return &Server{echo: e, container: c}
}

// Start begins listening on the configured port. Blocks until ctx is canceled.
func (s *Server) Start(ctx context.Context) error {
	// boilerplate:axis=rpc option=grpc START
	go func() {
		if err := s.container.GRPCServer.Start(); err != nil {
			log.Printf("grpc server error: %v", err)
		}
	}()
	// boilerplate:axis=rpc option=grpc END

	addr := fmt.Sprintf(":%d", s.container.Cfg.HTTP.Port)
	errCh := make(chan error, 1)
	go func() { errCh <- s.echo.Start(addr) }()
	select {
	case <-ctx.Done():
		// boilerplate:axis=rpc option=grpc START
		s.container.GRPCServer.Stop(ctx)
		// boilerplate:axis=rpc option=grpc END
		return s.echo.Shutdown(context.Background())
	case err := <-errCh:
		return err
	}
}

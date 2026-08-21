package app

import (
	"context"
	"fmt"
	"time"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/config"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/todo"

	// boilerplate:axis=cache option=redis START
	rediscache "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/cache/redis"
	// boilerplate:axis=cache option=redis END
	// boilerplate:axis=cache option=inmemory START
	inmemcache "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/cache/inmemory"
	// boilerplate:axis=cache option=inmemory END
	// boilerplate:axis=cache option=noop START
	noopcache "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/cache/noop"
	// boilerplate:axis=cache option=noop END
)

const defaultInMemoryCacheTTL = 5 * time.Minute

// newTodoCache builds the todo Cache for the selected cache backend.
// Mirrors newStorage: per-option marker-wrapped, unknown value fails fast.
func newTodoCache(ctx context.Context, cfg *config.Config) (todo.Cache, func(context.Context) error, error) {
	var valid []string
	// boilerplate:axis=cache option=redis START
	valid = append(valid, "redis")
	// boilerplate:axis=cache option=redis END
	// boilerplate:axis=cache option=inmemory START
	valid = append(valid, "inmemory")
	// boilerplate:axis=cache option=inmemory END
	// boilerplate:axis=cache option=noop START
	valid = append(valid, "noop")
	// boilerplate:axis=cache option=noop END

	switch cfg.CacheBackend {
	// boilerplate:axis=cache option=redis START
	case "redis":
		rdb, err := rediscache.Connect(ctx, cfg.Redis)
		if err != nil {
			return nil, nil, fmt.Errorf("app: redis: %w", err)
		}
		return rediscache.NewTodoCache(rdb), func(_ context.Context) error { return rdb.Close() }, nil
	// boilerplate:axis=cache option=redis END
	// boilerplate:axis=cache option=inmemory START
	case "inmemory":
		noClose := func(context.Context) error { return nil }
		return inmemcache.NewTodoCache(defaultInMemoryCacheTTL), noClose, nil
	// boilerplate:axis=cache option=inmemory END
	// boilerplate:axis=cache option=noop START
	case "noop":
		noClose := func(context.Context) error { return nil }
		return noopcache.NewTodoCache(), noClose, nil
	// boilerplate:axis=cache option=noop END
	}
	return nil, nil, fmt.Errorf("cache: CACHE_BACKEND must be one of %v; got %q", valid, cfg.CacheBackend)
}

package app

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/config"

	// boilerplate:axis=observability option=otel START
	otelobs "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/observability/otel"
	// boilerplate:axis=observability option=otel END
	// boilerplate:axis=observability option=datadog START
	datadogobs "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/observability/datadog"
	// boilerplate:axis=observability option=datadog END
	// boilerplate:axis=observability option=noop START
	noopobs "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/observability/noop"
	// boilerplate:axis=observability option=noop END
)

// observability bundles the telemetry wiring selected by
// OBSERVABILITY_BACKEND: the tracer shutdown func and the Echo middleware.
type observability struct {
	shutdown   func(context.Context) error
	middleware echo.MiddlewareFunc
}

// newObservability selects and initializes the observability adapter named
// by cfg.ObservabilityBackend. Each arm is marker-wrapped for atomic pruning
// (mirrors python/app/bootstrap/observability.py).
func newObservability(ctx context.Context, cfg *config.Config) (*observability, error) {
	switch cfg.ObservabilityBackend {
	// boilerplate:axis=observability option=otel START
	case "otel":
		shutdown, err := otelobs.Setup(ctx, cfg.OTel)
		if err != nil {
			return nil, fmt.Errorf("app: otel: %w", err)
		}
		return &observability{shutdown: shutdown, middleware: otelobs.Middleware(cfg.OTel.ServiceName)}, nil
	// boilerplate:axis=observability option=otel END
	// boilerplate:axis=observability option=datadog START
	case "datadog":
		shutdown, err := datadogobs.Setup(ctx, cfg.Datadog)
		if err != nil {
			return nil, fmt.Errorf("app: datadog: %w", err)
		}
		return &observability{shutdown: shutdown, middleware: datadogobs.Middleware(cfg.Datadog.ServiceName)}, nil
	// boilerplate:axis=observability option=datadog END
	// boilerplate:axis=observability option=noop START
	case "noop":
		shutdown, err := noopobs.Setup(ctx, noopobs.Config{})
		if err != nil {
			return nil, fmt.Errorf("app: noop: %w", err)
		}
		return &observability{shutdown: shutdown, middleware: noopobs.Middleware("")}, nil
	// boilerplate:axis=observability option=noop END
	default:
		return nil, fmt.Errorf("observability: unknown OBSERVABILITY_BACKEND %q", cfg.ObservabilityBackend)
	}
}

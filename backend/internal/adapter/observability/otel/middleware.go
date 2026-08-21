package otel

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

// Middleware returns an Echo middleware that auto-instruments incoming HTTP
// requests with spans.
func Middleware(serviceName string) echo.MiddlewareFunc {
	return otelecho.Middleware(serviceName)
}

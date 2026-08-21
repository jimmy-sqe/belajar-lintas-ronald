package datadog

import (
	"github.com/labstack/echo/v4"

	ddecho "gopkg.in/DataDog/dd-trace-go.v1/contrib/labstack/echo.v4"
)

// Middleware returns an Echo middleware that auto-instruments requests
// with Datadog APM spans.
func Middleware(serviceName string) echo.MiddlewareFunc {
	if serviceName == "" {
		serviceName = "backend-belajar-lintas-ronald"
	}
	return ddecho.Middleware(ddecho.WithServiceName(serviceName))
}

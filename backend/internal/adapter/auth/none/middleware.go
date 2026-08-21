// Package none is the no-auth adapter. It is a passthrough middleware that
// performs no authentication. Use for internal services or projects where
// auth happens at the edge (e.g., API gateway).
package none

import (
	"github.com/labstack/echo/v4"
)

// Middleware returns an Echo middleware that passes every request through
// without any auth check.
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return next
	}
}

package noop

import "github.com/labstack/echo/v4"

// Middleware returns an Echo middleware that does nothing.
func Middleware(_ string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return next
	}
}

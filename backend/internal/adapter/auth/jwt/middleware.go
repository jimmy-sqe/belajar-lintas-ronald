package jwt

import (
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/http/response"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/pkg/ctxutil"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/pkg/customerror"
)

// Middleware returns an Echo middleware that validates Bearer tokens via
// the given Verifier and injects user info into request context.
func Middleware(v *Verifier) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			h := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(h, "Bearer ") {
				return response.Err(c, customerror.ErrUnauthorized.WithMetadata(map[string]any{"detail": "Authorization header missing or malformed"}))
			}
			raw := strings.TrimPrefix(h, "Bearer ")
			claims, err := v.Parse(raw)
			if err != nil {
				return response.Err(c, customerror.ErrUnauthorized.WithMetadata(map[string]any{"detail": err.Error()}))
			}
			ctx := ctxutil.WithUserID(c.Request().Context(), claims.UserID)
			ctx = ctxutil.WithRoles(ctx, claims.Roles)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

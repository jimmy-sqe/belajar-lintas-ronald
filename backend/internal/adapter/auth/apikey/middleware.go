// Package apikey is an API-key auth adapter. It validates a fixed key
// from a configurable HTTP header against an allowlist.
package apikey

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/http/response"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/pkg/ctxutil"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/pkg/customerror"
)

// apiKeyNamespace derives a stable per-key principal UUID. Downstream handlers
// (e.g. the todo handler) parse the principal with uuid.Parse, so it MUST be a
// valid UUID — a raw "apikey:<prefix>" string would fail to parse and 401 every
// request. Mirrors Java's UUID.nameUUIDFromBytes("apikey:" + key).
var apiKeyNamespace = uuid.MustParse("6ba7b812-9dad-11d1-80b4-00c04fd430c8")

// Config holds API-key adapter settings.
type Config struct {
	HeaderName string   `mapstructure:"API_KEY_HEADER"`
	Keys       []string `mapstructure:"API_KEYS"` // comma-separated in env, parsed by viper
}

// Verifier checks API keys against the configured allowlist.
type Verifier struct{ cfg Config }

// NewVerifier returns a Verifier or an error if no keys are configured.
func NewVerifier(cfg Config) (*Verifier, error) {
	if cfg.HeaderName == "" {
		cfg.HeaderName = "X-API-Key"
	}
	if len(cfg.Keys) == 0 {
		return nil, errors.New("apikey: at least one key is required")
	}
	return &Verifier{cfg: cfg}, nil
}

// Verify returns (key, true) on match, ("", false) otherwise.
func (v *Verifier) Verify(presented string) (string, bool) {
	for _, k := range v.cfg.Keys {
		if strings.EqualFold(strings.TrimSpace(k), presented) {
			return k, true
		}
	}
	return "", false
}

// Middleware returns an Echo middleware that checks the API key header.
func Middleware(v *Verifier) echo.MiddlewareFunc {
	header := v.cfg.HeaderName
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			presented := c.Request().Header.Get(header)
			if presented == "" {
				return response.Err(c, customerror.ErrUnauthorized.WithMetadata(map[string]any{"detail": "API key header missing"}))
			}
			matched, ok := v.Verify(presented)
			if !ok {
				return response.Err(c, customerror.ErrUnauthorized.WithMetadata(map[string]any{"detail": "API key not recognised"}))
			}
			principal := uuid.NewSHA1(apiKeyNamespace, []byte("apikey:"+matched)).String()
			ctx := ctxutil.WithUserID(c.Request().Context(), principal)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

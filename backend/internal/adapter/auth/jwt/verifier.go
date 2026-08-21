// Package jwt is the JWT auth adapter. It validates Bearer tokens issued by
// an external IdP and injects the user identity into request context.
package jwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// Config holds JWT verification settings.
type Config struct {
	Secret        string `mapstructure:"JWT_SECRET"`
	Issuer        string `mapstructure:"JWT_ISSUER"`
	Audience      string `mapstructure:"JWT_AUDIENCE"`
	AccessTTLSec  int    `mapstructure:"JWT_ACCESS_TTL_SEC"`
	RefreshTTLSec int    `mapstructure:"JWT_REFRESH_TTL_SEC"`
}

// Claims is the standard set of claims this adapter expects.
type Claims struct {
	jwt.RegisteredClaims
	UserID string   `json:"sub"`
	Roles  []string `json:"roles,omitempty"`
}

// Verifier validates HS256 tokens against a shared secret.
type Verifier struct{ cfg Config }

// NewVerifier returns a Verifier or an error if Secret is empty.
func NewVerifier(cfg Config) (*Verifier, error) {
	if cfg.Secret == "" {
		return nil, errors.New("jwt: secret is required")
	}
	return &Verifier{cfg: cfg}, nil
}

// Parse validates the token and returns Claims.
func (v *Verifier) Parse(raw string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(raw, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt: unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(v.cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid {
		return nil, errors.New("jwt: invalid claims")
	}
	// Validate issuer if configured
	if v.cfg.Issuer != "" && claims.Issuer != v.cfg.Issuer {
		return nil, errors.New("jwt: invalid issuer")
	}
	// Validate audience if configured
	if v.cfg.Audience != "" {
		found := false
		for _, aud := range claims.Audience {
			if aud == v.cfg.Audience {
				found = true
				break
			}
		}
		if !found {
			return nil, errors.New("jwt: invalid audience")
		}
	}
	return claims, nil
}

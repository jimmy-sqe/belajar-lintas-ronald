// Package jwt provides an HS256 TokenIssuer implementation.
package jwt

import (
	"errors"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	authdomain "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/auth"
)

// Issuer signs and verifies tokens using HS256 and a shared secret.
type Issuer struct {
	secret []byte
}

// NewIssuer constructs an Issuer from a shared secret. Keep secret >= 32 chars.
func NewIssuer(secret string) *Issuer {
	return &Issuer{secret: []byte(secret)}
}

type customClaims struct {
	UserID string `json:"sub"`
	Email  string `json:"email"`
	Kind   string `json:"type"`
	gojwt.RegisteredClaims
}

// Sign produces a signed JWT string from the domain claims.
func (i *Issuer) Sign(c authdomain.Claims) (string, error) {
	tok := gojwt.NewWithClaims(gojwt.SigningMethodHS256, customClaims{
		UserID: c.UserID.String(),
		Email:  c.Email,
		Kind:   string(c.Kind),
		RegisteredClaims: gojwt.RegisteredClaims{
			IssuedAt:  gojwt.NewNumericDate(c.IssuedAt),
			ExpiresAt: gojwt.NewNumericDate(c.ExpiresAt),
		},
	})
	return tok.SignedString(i.secret)
}

// Parse verifies signature + expiry and returns the domain claims.
func (i *Issuer) Parse(token string) (authdomain.Claims, error) {
	var claims customClaims
	parsed, err := gojwt.ParseWithClaims(token, &claims, func(t *gojwt.Token) (any, error) {
		if t.Method.Alg() != gojwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return i.secret, nil
	})
	if err != nil {
		return authdomain.Claims{}, err
	}
	if !parsed.Valid {
		return authdomain.Claims{}, errors.New("invalid token")
	}

	id, err := uuid.Parse(claims.UserID)
	if err != nil {
		return authdomain.Claims{}, err
	}

	return authdomain.Claims{
		UserID:    id,
		Email:     claims.Email,
		Kind:      authdomain.TokenKind(claims.Kind),
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}

// Compile-time check that *Issuer satisfies authdomain.TokenIssuer.
var _ authdomain.TokenIssuer = (*Issuer)(nil)

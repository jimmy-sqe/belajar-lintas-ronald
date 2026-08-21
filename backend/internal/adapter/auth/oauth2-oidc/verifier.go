// Package oauth2oidc is an OAuth2/OIDC auth adapter. It validates Bearer
// tokens using the issuer's JWKS endpoint via go-oidc. Works with any
// OIDC-compliant IdP (Keycloak, Auth0, Cognito, Okta, Google, etc).
package oauth2oidc

import (
	"context"
	"errors"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
)

// Config holds OIDC verifier settings.
type Config struct {
	IssuerURL string `mapstructure:"OIDC_ISSUER_URL"`
	ClientID  string `mapstructure:"OIDC_CLIENT_ID"`
	Audience  string `mapstructure:"OIDC_AUDIENCE"`
}

// Verifier wraps an oidc.IDTokenVerifier.
type Verifier struct {
	v *oidc.IDTokenVerifier
}

// NewVerifier discovers the issuer and constructs a JWKS-backed verifier.
// Network call is made on construction; supply a context with timeout.
func NewVerifier(ctx context.Context, cfg Config) (*Verifier, error) {
	if cfg.IssuerURL == "" {
		return nil, errors.New("oauth2-oidc: issuer URL is required")
	}
	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("oauth2-oidc: discover provider: %w", err)
	}
	clientID := cfg.ClientID
	if clientID == "" {
		clientID = cfg.Audience
	}
	return &Verifier{v: provider.Verifier(&oidc.Config{ClientID: clientID})}, nil
}

// Parse verifies and parses the token, returning the OIDC ID token.
func (v *Verifier) Parse(ctx context.Context, raw string) (*oidc.IDToken, error) {
	tok, err := v.v.Verify(ctx, raw)
	if err != nil {
		return nil, fmt.Errorf("oauth2-oidc: verify: %w", err)
	}
	return tok, nil
}

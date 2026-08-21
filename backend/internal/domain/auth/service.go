package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SignInResult is what the service returns to handlers after a successful
// sign-in. Handlers shape this into the API envelope.
type SignInResult struct {
	AccessToken              string
	RefreshToken             string
	AccessTokenExpiresInSec  int
	RefreshTokenExpiresInSec int
	IssuedAt                 time.Time
	User                     *User
}

// RefreshResult is returned by Renew.
type RefreshResult struct {
	AccessToken             string
	AccessTokenExpiresInSec int
	IssuedAt                time.Time
}

// Authenticator orchestrates login / logout / renew flows.
type Authenticator interface {
	SignIn(ctx context.Context, email, password string) (*SignInResult, error)
	Logout(ctx context.Context, refreshToken string) error
	Renew(ctx context.Context, refreshToken string) (*RefreshResult, error)
}

// ServiceConfig captures TTLs read from app config so the service is
// pure (no env lookup inside).
type ServiceConfig struct {
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type service struct {
	users    UserRepository
	sessions SessionStore
	tokens   TokenIssuer
	cfg      ServiceConfig
}

// NewService constructs an Authenticator with explicit dependencies.
func NewService(users UserRepository, sessions SessionStore,
	tokens TokenIssuer, cfg ServiceConfig) Authenticator {
	return &service{users: users, sessions: sessions, tokens: tokens, cfg: cfg}
}

func (s *service) SignIn(ctx context.Context, email, password string) (*SignInResult, error) {
	u, err := s.users.FindByEmail(ctx, strings.ToLower(email))
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	now := time.Now().UTC()
	access, err := s.tokens.Sign(Claims{
		UserID: u.ID, Email: u.Email, Kind: KindAccess,
		IssuedAt: now, ExpiresAt: now.Add(s.cfg.AccessTTL),
	})
	if err != nil {
		return nil, err
	}
	refresh, err := s.tokens.Sign(Claims{
		UserID: u.ID, Email: u.Email, Kind: KindRefresh,
		IssuedAt: now, ExpiresAt: now.Add(s.cfg.RefreshTTL),
	})
	if err != nil {
		return nil, err
	}

	if err := s.sessions.Put(ctx, RefreshSession{
		TokenHash: hashToken(refresh),
		UserID:    u.ID,
		ExpiresAt: now.Add(s.cfg.RefreshTTL),
	}); err != nil {
		return nil, err
	}

	return &SignInResult{
		AccessToken:              access,
		RefreshToken:             refresh,
		AccessTokenExpiresInSec:  int(s.cfg.AccessTTL.Seconds()),
		RefreshTokenExpiresInSec: int(s.cfg.RefreshTTL.Seconds()),
		IssuedAt:                 now,
		User:                     u,
	}, nil
}

func (s *service) Logout(ctx context.Context, refreshToken string) error {
	return s.sessions.Delete(ctx, hashToken(refreshToken))
}

func (s *service) Renew(ctx context.Context, refreshToken string) (*RefreshResult, error) {
	claims, err := s.tokens.Parse(refreshToken)
	if err != nil {
		return nil, ErrRefreshInvalid
	}
	if claims.Kind != KindRefresh {
		return nil, ErrRefreshInvalid
	}
	rec, err := s.sessions.Get(ctx, hashToken(refreshToken))
	if err != nil || rec == nil {
		return nil, ErrRefreshInvalid
	}
	if rec.UserID != claims.UserID {
		return nil, ErrRefreshInvalid
	}
	if time.Now().UTC().After(rec.ExpiresAt) {
		_ = s.sessions.Delete(ctx, hashToken(refreshToken))
		return nil, ErrRefreshInvalid
	}

	now := time.Now().UTC()
	access, err := s.tokens.Sign(Claims{
		UserID: claims.UserID, Email: claims.Email, Kind: KindAccess,
		IssuedAt: now, ExpiresAt: now.Add(s.cfg.AccessTTL),
	})
	if err != nil {
		return nil, err
	}

	return &RefreshResult{
		AccessToken:             access,
		AccessTokenExpiresInSec: int(s.cfg.AccessTTL.Seconds()),
		IssuedAt:                now,
	}, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// RefreshSession represents a server-side refresh-token record stored
// in Redis. The token itself is hashed before being used as a key.
type RefreshSession struct {
	TokenHash string
	UserID    uuid.UUID
	ExpiresAt time.Time
}

// SessionStore abstracts the refresh-token storage backend (Redis).
type SessionStore interface {
	Put(ctx context.Context, session RefreshSession) error
	Get(ctx context.Context, tokenHash string) (*RefreshSession, error)
	Delete(ctx context.Context, tokenHash string) error
}

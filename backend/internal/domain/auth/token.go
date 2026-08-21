package auth

import (
	"time"

	"github.com/google/uuid"
)

// TokenKind discriminates JWT token types so middleware can reject
// refresh tokens presented in the Authorization header.
type TokenKind string

const (
	KindAccess  TokenKind = "access"
	KindRefresh TokenKind = "refresh"
)

// Claims carries the JWT payload used by both access and refresh tokens.
type Claims struct {
	UserID    uuid.UUID
	Email     string
	Kind      TokenKind
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// TokenIssuer signs and verifies JWT tokens. Implementations live in
// internal/adapter/token/<algo>/.
type TokenIssuer interface {
	Sign(claims Claims) (string, error)
	Parse(token string) (Claims, error)
}

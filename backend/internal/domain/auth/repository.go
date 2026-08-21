package auth

import (
	"context"

	"github.com/google/uuid"
)

// UserRepository abstracts user storage across persistence dialects
// (postgres, mysql, mongodb).
type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
}

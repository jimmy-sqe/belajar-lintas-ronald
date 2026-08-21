// Package auth handles user authentication and session management.
package auth

import (
	"time"

	"github.com/google/uuid"
)

// User is the canonical authenticated identity. PasswordHash is bcrypt
// and must never appear in API DTOs or logs.
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Name         string
	CreatedAt    time.Time
	CreatedBy    *uuid.UUID
	ModifiedAt   time.Time
	ModifiedBy   *uuid.UUID
}

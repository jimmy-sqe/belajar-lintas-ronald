//go:build integration
// +build integration

package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pgadapter "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/persistence/postgres"
	authdomain "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/auth"
)

func TestUserRepository_FindByEmail_FindByID(t *testing.T) {
	// Reuse existing test infrastructure: setupDB spins up a testcontainer
	// and runs all migrations (including 001_create_users from Task 6).
	ctx := context.Background()
	db, cleanup := setupDB(t)
	defer cleanup()

	id := uuid.New()
	_, err := db.ExecContext(ctx, `INSERT INTO users (id, email, password_hash, name)
		VALUES ($1, $2, $3, $4)`,
		id.String(), "demo@example.com", "$2a$12$bcrypt-placeholder-hash", "Demo")
	require.NoError(t, err)

	repo := pgadapter.NewUserRepository(db)

	t.Run("FindByEmail success", func(t *testing.T) {
		u, err := repo.FindByEmail(ctx, "demo@example.com")
		require.NoError(t, err)
		assert.Equal(t, id, u.ID)
		assert.Equal(t, "demo@example.com", u.Email)
	})

	t.Run("FindByEmail case-insensitive", func(t *testing.T) {
		u, err := repo.FindByEmail(ctx, "DEMO@example.com")
		require.NoError(t, err)
		assert.Equal(t, id, u.ID)
	})

	t.Run("FindByEmail not found", func(t *testing.T) {
		_, err := repo.FindByEmail(ctx, "nobody@example.com")
		assert.ErrorIs(t, err, authdomain.ErrUserNotFound)
	})

	t.Run("FindByID success", func(t *testing.T) {
		u, err := repo.FindByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "demo@example.com", u.Email)
	})

	t.Run("FindByID not found", func(t *testing.T) {
		_, err := repo.FindByID(ctx, uuid.New())
		assert.ErrorIs(t, err, authdomain.ErrUserNotFound)
	})
}

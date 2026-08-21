//go:build integration

package mongodb_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcmongo "github.com/testcontainers/testcontainers-go/modules/mongodb"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/persistence/mongodb"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/todo"
)

func setup(t *testing.T) (*mongodb.TodoRepository, func()) {
	t.Helper()
	ctx := context.Background()
	c, err := tcmongo.RunContainer(ctx, testcontainers.WithImage("mongo:7"))
	require.NoError(t, err)
	uri, err := c.ConnectionString(ctx)
	require.NoError(t, err)

	client, err := mongodb.Connect(ctx, mongodb.Config{URI: uri, DB: "test"})
	require.NoError(t, err)

	repo := mongodb.NewTodoRepository(client.Database("test"))
	return repo, func() {
		_ = client.Disconnect(context.Background())
		_ = c.Terminate(ctx)
	}
}

func TestTodoRepository_CreateAndFindByID(t *testing.T) {
	repo, cleanup := setup(t)
	defer cleanup()

	ctx := context.Background()
	ownerID := uuid.New()
	now := time.Now().UTC().Truncate(time.Millisecond)
	in := &todo.Todo{
		ID:         uuid.New(),
		Title:      "Test Todo",
		CreatedAt:  now,
		CreatedBy:  ownerID,
		ModifiedAt: now,
		ModifiedBy: ownerID,
	}
	require.NoError(t, repo.Create(ctx, in))
	got, err := repo.FindByID(ctx, in.ID)
	require.NoError(t, err)
	require.Equal(t, in.Title, got.Title)
}

func TestTodoRepository_FindByID_NotFound(t *testing.T) {
	repo, cleanup := setup(t)
	defer cleanup()
	got, err := repo.FindByID(context.Background(), uuid.New())
	require.ErrorIs(t, err, todo.ErrNotFound)
	require.Nil(t, got)
}

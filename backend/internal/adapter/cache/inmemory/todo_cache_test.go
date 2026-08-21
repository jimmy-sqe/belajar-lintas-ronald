package inmemory_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/cache/inmemory"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/todo"
)

func TestTodoCache_SetThenGet(t *testing.T) {
	c := inmemory.NewTodoCache(time.Minute)
	ctx := context.Background()
	id := uuid.New()
	in := &todo.Todo{ID: id, Title: "Foo"}

	require.NoError(t, c.Set(ctx, in))
	got, ok := c.Get(ctx, id)
	require.True(t, ok)
	require.NotNil(t, got)
	require.Equal(t, "Foo", got.Title)
}

func TestTodoCache_GetMiss(t *testing.T) {
	c := inmemory.NewTodoCache(time.Minute)
	got, ok := c.Get(context.Background(), uuid.New())
	require.False(t, ok)
	require.Nil(t, got)
}

func TestTodoCache_TTLExpiry(t *testing.T) {
	c := inmemory.NewTodoCache(20 * time.Millisecond)
	ctx := context.Background()
	id := uuid.New()
	require.NoError(t, c.Set(ctx, &todo.Todo{ID: id, Title: "Foo"}))

	time.Sleep(40 * time.Millisecond)
	got, ok := c.Get(ctx, id)
	require.False(t, ok)
	require.Nil(t, got)
}

func TestTodoCache_Delete(t *testing.T) {
	c := inmemory.NewTodoCache(time.Minute)
	ctx := context.Background()
	id := uuid.New()
	require.NoError(t, c.Set(ctx, &todo.Todo{ID: id}))
	require.NoError(t, c.Delete(ctx, id))
	got, ok := c.Get(ctx, id)
	require.False(t, ok)
	require.Nil(t, got)
}

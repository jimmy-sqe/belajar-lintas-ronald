package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/todo"
)

// Cache is a mock implementation of todo.Cache.
type Cache struct{ mock.Mock }

func (m *Cache) Get(ctx context.Context, id uuid.UUID) (*todo.Todo, bool) {
	args := m.Called(ctx, id)
	v, _ := args.Get(0).(*todo.Todo)
	return v, args.Bool(1)
}

func (m *Cache) Set(ctx context.Context, t *todo.Todo) error {
	return m.Called(ctx, t).Error(0)
}

func (m *Cache) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

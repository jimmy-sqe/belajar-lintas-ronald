// Package mock holds handwritten mocks for the todo ports. They use
// stretchr/testify/mock.
package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/todo"
)

// Repository is a mock implementation of todo.Repository.
type Repository struct{ mock.Mock }

func (m *Repository) Create(ctx context.Context, t *todo.Todo) error {
	return m.Called(ctx, t).Error(0)
}

func (m *Repository) FindByID(ctx context.Context, id uuid.UUID) (*todo.Todo, error) {
	args := m.Called(ctx, id)
	v, _ := args.Get(0).(*todo.Todo)
	return v, args.Error(1)
}

func (m *Repository) List(ctx context.Context, filter todo.ListFilter) ([]todo.Todo, int64, error) {
	args := m.Called(ctx, filter)
	v, _ := args.Get(0).([]todo.Todo)
	return v, args.Get(1).(int64), args.Error(2)
}

func (m *Repository) Update(ctx context.Context, t *todo.Todo) error {
	return m.Called(ctx, t).Error(0)
}

func (m *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

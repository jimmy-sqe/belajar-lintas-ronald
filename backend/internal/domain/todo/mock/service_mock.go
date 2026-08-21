package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/todo"
)

// Service is a mock implementation of todo.Service for handler tests.
type Service struct{ mock.Mock }

func (m *Service) Create(ctx context.Context, input todo.CreateInput) (*todo.Todo, error) {
	args := m.Called(ctx, input)
	v, _ := args.Get(0).(*todo.Todo)
	return v, args.Error(1)
}

func (m *Service) FindByID(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) (*todo.Todo, error) {
	args := m.Called(ctx, id, ownerID)
	v, _ := args.Get(0).(*todo.Todo)
	return v, args.Error(1)
}

func (m *Service) List(ctx context.Context, filter todo.ListFilter) ([]todo.Todo, int64, error) {
	args := m.Called(ctx, filter)
	v, _ := args.Get(0).([]todo.Todo)
	return v, args.Get(1).(int64), args.Error(2)
}

func (m *Service) Update(ctx context.Context, id uuid.UUID, input todo.UpdateInput) (*todo.Todo, error) {
	args := m.Called(ctx, id, input)
	v, _ := args.Get(0).(*todo.Todo)
	return v, args.Error(1)
}

func (m *Service) Delete(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) error {
	return m.Called(ctx, id, ownerID).Error(0)
}

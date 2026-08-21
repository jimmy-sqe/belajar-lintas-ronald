package todo

import (
	"context"

	"github.com/google/uuid"
)

// Repository persists todos. Implementations live in
// internal/adapter/persistence/<dialect>/.
type Repository interface {
	Create(ctx context.Context, todo *Todo) error
	FindByID(ctx context.Context, id uuid.UUID) (*Todo, error)
	List(ctx context.Context, filter ListFilter) ([]Todo, int64, error)
	Update(ctx context.Context, todo *Todo) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// Cache fronts Repository for hot reads. Implementations live in
// internal/adapter/cache/<impl>/.
type Cache interface {
	Get(ctx context.Context, id uuid.UUID) (*Todo, bool)
	Set(ctx context.Context, todo *Todo) error
	Delete(ctx context.Context, id uuid.UUID) error
}

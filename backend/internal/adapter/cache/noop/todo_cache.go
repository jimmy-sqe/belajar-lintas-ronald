// Package noop is the no-op cache adapter. It implements todo.Cache without
// any actual caching — useful for projects where caching adds complexity
// without enough benefit.
package noop

import (
	"context"

	"github.com/google/uuid"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/todo"
)

// TodoCache is the no-op implementation of todo.Cache.
type TodoCache struct{}

// NewTodoCache constructs a no-op TodoCache.
func NewTodoCache() *TodoCache { return &TodoCache{} }

// Get always returns (nil, false).
func (c *TodoCache) Get(_ context.Context, _ uuid.UUID) (*todo.Todo, bool) { return nil, false }

// Set is a no-op that returns nil.
func (c *TodoCache) Set(_ context.Context, _ *todo.Todo) error { return nil }

// Delete is a no-op that returns nil.
func (c *TodoCache) Delete(_ context.Context, _ uuid.UUID) error { return nil }

var _ todo.Cache = (*TodoCache)(nil)

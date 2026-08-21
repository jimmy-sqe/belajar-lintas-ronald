// Package inmemory is the in-memory cache adapter. It implements todo.Cache
// using sync.Map with per-entry TTL. Suitable for single-instance services
// or caches that don't need to be shared across pods.
package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/todo"
)

// TodoCache is the in-memory implementation of todo.Cache.
type TodoCache struct {
	store sync.Map // map[string]entry
	ttl   time.Duration
}

type entry struct {
	todo      *todo.Todo
	expiresAt time.Time
}

// NewTodoCache constructs a TodoCache with the given TTL.
func NewTodoCache(ttl time.Duration) *TodoCache {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &TodoCache{ttl: ttl}
}

// Get returns (nil, false) on cache miss or after TTL expiry; (*Todo, true) on hit.
func (c *TodoCache) Get(_ context.Context, id uuid.UUID) (*todo.Todo, bool) {
	v, ok := c.store.Load(id.String())
	if !ok {
		return nil, false
	}
	e, _ := v.(entry)
	if time.Now().After(e.expiresAt) {
		c.store.Delete(id.String())
		return nil, false
	}
	cp := *e.todo
	return &cp, true
}

// Set stores the todo with the configured TTL.
func (c *TodoCache) Set(_ context.Context, t *todo.Todo) error {
	cp := *t
	c.store.Store(t.ID.String(), entry{todo: &cp, expiresAt: time.Now().Add(c.ttl)})
	return nil
}

// Delete removes the cache entry. Missing keys are not an error.
func (c *TodoCache) Delete(_ context.Context, id uuid.UUID) error {
	c.store.Delete(id.String())
	return nil
}

var _ todo.Cache = (*TodoCache)(nil)

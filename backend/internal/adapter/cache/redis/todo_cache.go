package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/todo"
)

const (
	todoKeyPrefix = "todo:"
	todoTTL       = 5 * time.Minute
)

// TodoCache is the Redis implementation of todo.Cache.
type TodoCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewTodoCache constructs a TodoCache with default 5-minute TTL.
func NewTodoCache(client *redis.Client) *TodoCache {
	return &TodoCache{client: client, ttl: todoTTL}
}

func key(id uuid.UUID) string { return todoKeyPrefix + id.String() }

// Get returns (nil, false) on cache miss; (*Todo, true) on hit.
func (c *TodoCache) Get(ctx context.Context, id uuid.UUID) (*todo.Todo, bool) {
	raw, err := c.client.Get(ctx, key(id)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false
		}
		return nil, false
	}
	var t todo.Todo
	if err := json.Unmarshal(raw, &t); err != nil {
		// corrupt cache entry: treat as miss and invalidate
		_ = c.client.Del(ctx, key(id)).Err()
		return nil, false
	}
	return &t, true
}

// Set stores the todo with the configured TTL.
func (c *TodoCache) Set(ctx context.Context, t *todo.Todo) error {
	raw, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("redis: marshal: %w", err)
	}
	if err := c.client.Set(ctx, key(t.ID), raw, c.ttl).Err(); err != nil {
		return fmt.Errorf("redis: set: %w", err)
	}
	return nil
}

// Delete removes the cache entry. Missing keys are not an error.
func (c *TodoCache) Delete(ctx context.Context, id uuid.UUID) error {
	if err := c.client.Del(ctx, key(id)).Err(); err != nil {
		return fmt.Errorf("redis: del: %w", err)
	}
	return nil
}

var _ todo.Cache = (*TodoCache)(nil)

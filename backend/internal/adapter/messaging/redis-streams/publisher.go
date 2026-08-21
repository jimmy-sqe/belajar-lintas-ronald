package redisstreams

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// Publish appends an entry to the given stream. `topic` is the stream key.
func (p *Publisher) Publish(ctx context.Context, topic string, payload []byte) error {
	if err := p.client.XAdd(ctx, &redis.XAddArgs{
		Stream: topic,
		Values: map[string]any{"payload": payload},
	}).Err(); err != nil {
		return fmt.Errorf("redis-streams: xadd: %w", err)
	}
	return nil
}

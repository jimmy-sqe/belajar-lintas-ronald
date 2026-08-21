// Package redisstreams is the Redis Streams messaging adapter.
package redisstreams

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Config holds Redis Streams settings (shares Redis infra with cache axis).
type Config struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     int    `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
}

// Publisher wraps a redis.Client.
type Publisher struct {
	client *redis.Client
}

// New connects to Redis.
func New(ctx context.Context, cfg Config) (*Publisher, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
	})
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := c.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("redis-streams: ping: %w", err)
	}
	return &Publisher{client: c}, nil
}

// Close releases the underlying client.
func (p *Publisher) Close() error { return p.client.Close() }

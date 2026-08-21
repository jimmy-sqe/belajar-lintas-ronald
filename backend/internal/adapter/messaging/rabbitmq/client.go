// Package rabbitmq is the RabbitMQ messaging adapter via amqp091.
package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Config holds RabbitMQ settings.
type Config struct {
	URL string `mapstructure:"RABBITMQ_URL"`
}

// Publisher wraps an amqp connection + channel.
type Publisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// New connects to RabbitMQ and opens a publisher channel.
func New(_ context.Context, cfg Config) (*Publisher, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("rabbitmq: URL is required")
	}
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq: dial: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("rabbitmq: channel: %w", err)
	}
	return &Publisher{conn: conn, ch: ch}, nil
}

// Close shuts down the channel and connection.
func (p *Publisher) Close() error {
	_ = p.ch.Close()
	return p.conn.Close()
}

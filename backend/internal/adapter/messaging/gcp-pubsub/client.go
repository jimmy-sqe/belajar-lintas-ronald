// Package gcppubsub is the GCP Pub/Sub messaging adapter.
package gcppubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

// Config holds GCP Pub/Sub settings.
type Config struct {
	ProjectID string `mapstructure:"GCP_PROJECT_ID"`
	Topic     string `mapstructure:"PUBSUB_TOPIC"`
}

// Publisher wraps a pubsub.Client + a topic name.
type Publisher struct {
	client *pubsub.Client
	topic  string
}

// New constructs a Publisher.
func New(ctx context.Context, cfg Config) (*Publisher, error) {
	if cfg.ProjectID == "" {
		return nil, fmt.Errorf("gcp-pubsub: project ID is required")
	}
	c, err := pubsub.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("gcp-pubsub: client: %w", err)
	}
	return &Publisher{client: c, topic: cfg.Topic}, nil
}

// Close releases the underlying client.
func (p *Publisher) Close() error { return p.client.Close() }

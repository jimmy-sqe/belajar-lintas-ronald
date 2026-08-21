package gcppubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

// Publish sends payload to the given topic. If `topic` is empty, falls back
// to the default topic from Config.
func (p *Publisher) Publish(ctx context.Context, topic string, payload []byte) error {
	if topic == "" {
		topic = p.topic
	}
	if topic == "" {
		return fmt.Errorf("gcp-pubsub: topic is empty")
	}
	t := p.client.Topic(topic)
	res := t.Publish(ctx, &pubsub.Message{Data: payload})
	if _, err := res.Get(ctx); err != nil {
		return fmt.Errorf("gcp-pubsub: publish: %w", err)
	}
	return nil
}

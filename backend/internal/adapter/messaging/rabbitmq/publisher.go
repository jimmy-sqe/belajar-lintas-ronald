package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publish sends payload to the named exchange (using empty routing key).
// Treats `topic` as the exchange name.
func (p *Publisher) Publish(ctx context.Context, topic string, payload []byte) error {
	if err := p.ch.PublishWithContext(ctx, topic, "", false, false, amqp.Publishing{
		ContentType: "application/octet-stream",
		Body:        payload,
	}); err != nil {
		return fmt.Errorf("rabbitmq: publish: %w", err)
	}
	return nil
}

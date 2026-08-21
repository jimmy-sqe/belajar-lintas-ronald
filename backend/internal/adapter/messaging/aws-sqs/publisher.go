package awssqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// Publish sends a message to the SQS queue. `topic` is the queue URL; if
// empty, falls back to default from Config.
func (p *Publisher) Publish(ctx context.Context, topic string, payload []byte) error {
	url := topic
	if url == "" {
		url = p.queueURL
	}
	if url == "" {
		return fmt.Errorf("aws-sqs: queue URL is empty")
	}
	body := string(payload)
	if _, err := p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &url,
		MessageBody: &body,
	}); err != nil {
		return fmt.Errorf("aws-sqs: send: %w", err)
	}
	return nil
}

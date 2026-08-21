package awssns

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// Publish sends a message to the SNS topic. `topic` is treated as the
// topic ARN; if empty, falls back to default from Config.
func (p *Publisher) Publish(ctx context.Context, topic string, payload []byte) error {
	arn := topic
	if arn == "" {
		arn = p.topicARN
	}
	if arn == "" {
		return fmt.Errorf("aws-sns: topic ARN is empty")
	}
	msg := string(payload)
	if _, err := p.client.Publish(ctx, &sns.PublishInput{
		TopicArn: &arn,
		Message:  &msg,
	}); err != nil {
		return fmt.Errorf("aws-sns: publish: %w", err)
	}
	return nil
}

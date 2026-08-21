// Package awssqs is the AWS SQS messaging adapter.
package awssqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// Config holds AWS SQS settings.
type Config struct {
	Region   string `mapstructure:"AWS_REGION"`
	QueueURL string `mapstructure:"SQS_QUEUE_URL"`
}

// Publisher wraps the SQS client + default queue URL.
type Publisher struct {
	client   *sqs.Client
	queueURL string
}

// New constructs a Publisher.
func New(ctx context.Context, cfg Config) (*Publisher, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("aws-sqs: load aws config: %w", err)
	}
	return &Publisher{client: sqs.NewFromConfig(awsCfg), queueURL: cfg.QueueURL}, nil
}

// Close is a no-op for SQS.
func (p *Publisher) Close() error { return nil }

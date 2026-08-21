// Package awssns is the AWS SNS messaging adapter.
package awssns

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// Config holds AWS SNS settings.
type Config struct {
	Region   string `mapstructure:"AWS_REGION"`
	TopicARN string `mapstructure:"SNS_TOPIC_ARN"`
}

// Publisher wraps the SNS client + default topic ARN.
type Publisher struct {
	client   *sns.Client
	topicARN string
}

// New constructs a Publisher.
func New(ctx context.Context, cfg Config) (*Publisher, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("aws-sns: load aws config: %w", err)
	}
	return &Publisher{client: sns.NewFromConfig(awsCfg), topicARN: cfg.TopicARN}, nil
}

// Close is a no-op for SNS.
func (p *Publisher) Close() error { return nil }

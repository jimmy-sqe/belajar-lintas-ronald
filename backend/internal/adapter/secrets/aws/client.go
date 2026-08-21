// Package aws is the AWS Secrets Manager adapter.
package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// Config holds AWS Secrets Manager settings.
type Config struct {
	Region string `mapstructure:"AWS_REGION"`
}

// Client wraps the Secrets Manager client.
type Client struct {
	client *secretsmanager.Client
}

// New constructs a Client.
func New(ctx context.Context, cfg Config) (*Client, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("secrets/aws: load aws config: %w", err)
	}
	return &Client{client: secretsmanager.NewFromConfig(awsCfg)}, nil
}

// Get retrieves the named secret. Returns the SecretString content.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	out, err := c.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{SecretId: &key})
	if err != nil {
		return "", fmt.Errorf("secrets/aws: get: %w", err)
	}
	if out.SecretString == nil {
		return "", fmt.Errorf("secrets/aws: %s has no SecretString", key)
	}
	return *out.SecretString, nil
}

// Package s3 is the AWS S3 storage adapter.
package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

// Config holds S3 settings.
type Config struct {
	Region string `mapstructure:"AWS_REGION"`
	Bucket string `mapstructure:"S3_BUCKET"`
}

// Client wraps an S3 client + bucket.
type Client struct {
	client *awss3.Client
	bucket string
}

// New constructs a Client.
func New(ctx context.Context, cfg Config) (*Client, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("s3: bucket is required")
	}
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("s3: load aws config: %w", err)
	}
	return &Client{client: awss3.NewFromConfig(awsCfg), bucket: cfg.Bucket}, nil
}

func (c *Client) Upload(ctx context.Context, key string, content io.Reader) error {
	if _, err := c.client.PutObject(ctx, &awss3.PutObjectInput{
		Bucket: &c.bucket,
		Key:    &key,
		Body:   content,
	}); err != nil {
		return fmt.Errorf("s3: put: %w", err)
	}
	return nil
}

func (c *Client) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	out, err := c.client.GetObject(ctx, &awss3.GetObjectInput{
		Bucket: &c.bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, fmt.Errorf("s3: get: %w", err)
	}
	return out.Body, nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	if _, err := c.client.DeleteObject(ctx, &awss3.DeleteObjectInput{
		Bucket: &c.bucket,
		Key:    &key,
	}); err != nil {
		return fmt.Errorf("s3: delete: %w", err)
	}
	return nil
}

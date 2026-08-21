// Package gcs is the Google Cloud Storage adapter.
package gcs

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

// Config holds GCS settings.
type Config struct {
	Bucket string `mapstructure:"GCS_BUCKET"`
}

// Client wraps a storage.Client + a default bucket.
type Client struct {
	client *storage.Client
	bucket string
}

// New constructs a Client.
func New(ctx context.Context, cfg Config) (*Client, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("gcs: bucket is required")
	}
	c, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("gcs: client: %w", err)
	}
	return &Client{client: c, bucket: cfg.Bucket}, nil
}

func (c *Client) Upload(ctx context.Context, key string, content io.Reader) error {
	w := c.client.Bucket(c.bucket).Object(key).NewWriter(ctx)
	if _, err := io.Copy(w, content); err != nil {
		_ = w.Close()
		return fmt.Errorf("gcs: copy: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("gcs: close writer: %w", err)
	}
	return nil
}

func (c *Client) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	r, err := c.client.Bucket(c.bucket).Object(key).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("gcs: reader: %w", err)
	}
	return r, nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	if err := c.client.Bucket(c.bucket).Object(key).Delete(ctx); err != nil {
		return fmt.Errorf("gcs: delete: %w", err)
	}
	return nil
}

// Close releases the underlying client.
func (c *Client) Close() error { return c.client.Close() }

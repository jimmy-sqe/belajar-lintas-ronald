// Package minio is the MinIO storage adapter (S3-compatible self-hosted).
package minio

import (
	"context"
	"fmt"
	"io"

	mc "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Config holds MinIO settings.
type Config struct {
	Endpoint  string `mapstructure:"MINIO_ENDPOINT"`
	AccessKey string `mapstructure:"MINIO_ACCESS_KEY"`
	SecretKey string `mapstructure:"MINIO_SECRET_KEY"`
	Bucket    string `mapstructure:"MINIO_BUCKET"`
	UseSSL    bool   `mapstructure:"MINIO_USE_SSL"`
}

// Client wraps a MinIO client + bucket.
type Client struct {
	client *mc.Client
	bucket string
}

// New constructs a Client and verifies the bucket exists (or creates it).
func New(ctx context.Context, cfg Config) (*Client, error) {
	if cfg.Endpoint == "" || cfg.Bucket == "" {
		return nil, fmt.Errorf("minio: endpoint and bucket are required")
	}
	c, err := mc.New(cfg.Endpoint, &mc.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio: new: %w", err)
	}
	exists, err := c.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("minio: bucket exists: %w", err)
	}
	if !exists {
		if err := c.MakeBucket(ctx, cfg.Bucket, mc.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("minio: make bucket: %w", err)
		}
	}
	return &Client{client: c, bucket: cfg.Bucket}, nil
}

func (c *Client) Upload(ctx context.Context, key string, content io.Reader) error {
	if _, err := c.client.PutObject(ctx, c.bucket, key, content, -1, mc.PutObjectOptions{}); err != nil {
		return fmt.Errorf("minio: put: %w", err)
	}
	return nil
}

func (c *Client) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	obj, err := c.client.GetObject(ctx, c.bucket, key, mc.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("minio: get: %w", err)
	}
	return obj, nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	if err := c.client.RemoveObject(ctx, c.bucket, key, mc.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("minio: remove: %w", err)
	}
	return nil
}

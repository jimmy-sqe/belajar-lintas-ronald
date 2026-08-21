// Package gcp is the Google Secret Manager adapter.
package gcp

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// Config holds GCP Secret Manager settings.
type Config struct {
	ProjectID string `mapstructure:"GCP_PROJECT_ID"`
}

// Client wraps a Secret Manager client.
type Client struct {
	client    *secretmanager.Client
	projectID string
}

// New constructs a Client.
func New(ctx context.Context, cfg Config) (*Client, error) {
	if cfg.ProjectID == "" {
		return nil, fmt.Errorf("secrets/gcp: project ID is required")
	}
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("secrets/gcp: client: %w", err)
	}
	return &Client{client: c, projectID: cfg.ProjectID}, nil
}

// Get retrieves the latest version of the secret named by key.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", c.projectID, key)
	res, err := c.client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{Name: name})
	if err != nil {
		return "", fmt.Errorf("secrets/gcp: access: %w", err)
	}
	return string(res.Payload.Data), nil
}

// Close releases the underlying client.
func (c *Client) Close() error { return c.client.Close() }

// Package local is the filesystem-based storage adapter.
package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Config holds local storage settings.
type Config struct {
	BasePath string `mapstructure:"LOCAL_STORAGE_PATH"`
}

// Client implements storage.Storage using the local filesystem.
type Client struct {
	base string
}

// New returns a Client. The base path is created if it doesn't exist.
func New(cfg Config) (*Client, error) {
	if cfg.BasePath == "" {
		return nil, fmt.Errorf("local: BasePath is required")
	}
	if err := os.MkdirAll(cfg.BasePath, 0o755); err != nil { //nolint:gosec // 0o755 is required so other processes (e.g. nginx, sidecars) can read uploaded assets
		return nil, fmt.Errorf("local: mkdir base: %w", err)
	}
	return &Client{base: cfg.BasePath}, nil
}

func (c *Client) path(key string) string { return filepath.Join(c.base, key) }

// Upload writes the reader contents to the keyed path.
func (c *Client) Upload(_ context.Context, key string, content io.Reader) error {
	p := c.path(key)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil { //nolint:gosec // 0o755 is required so other processes (e.g. nginx, sidecars) can read uploaded assets
		return fmt.Errorf("local: mkdir parent: %w", err)
	}
	f, err := os.Create(p) //nolint:gosec // p is rooted at the validated cfg.BasePath; key traversal is the caller's responsibility (domain layer sanitises)
	if err != nil {
		return fmt.Errorf("local: create: %w", err)
	}
	defer f.Close() //nolint:errcheck // best-effort close; write errors propagate via io.Copy
	if _, err := io.Copy(f, content); err != nil {
		return fmt.Errorf("local: copy: %w", err)
	}
	return nil
}

// Download opens the file at the keyed path.
func (c *Client) Download(_ context.Context, key string) (io.ReadCloser, error) {
	f, err := os.Open(c.path(key))
	if err != nil {
		return nil, fmt.Errorf("local: open: %w", err)
	}
	return f, nil
}

// Delete removes the file at the keyed path. Missing file is not an error.
func (c *Client) Delete(_ context.Context, key string) error {
	if err := os.Remove(c.path(key)); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("local: remove: %w", err)
	}
	return nil
}

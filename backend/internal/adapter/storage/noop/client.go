// Package noop is the no-op storage adapter. Upload/Delete are silent successes.
// Download returns an empty reader.
package noop

import (
	"bytes"
	"context"
	"io"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/storage"
)

// Client is the no-op implementation.
type Client struct{}

// New returns a no-op Client.
func New() storage.Storage { return &Client{} }

// Upload silently succeeds.
func (c *Client) Upload(_ context.Context, _ string, _ io.Reader) error { return nil }

// Download returns an empty reader.
func (c *Client) Download(_ context.Context, _ string) (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader(nil)), nil
}

// Delete silently succeeds.
func (c *Client) Delete(_ context.Context, _ string) error { return nil }

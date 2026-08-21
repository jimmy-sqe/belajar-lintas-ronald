// Package storage defines the common Storage port. All option subpackages
// implement this interface.
package storage

import (
	"context"
	"io"
)

// Storage uploads, downloads, and deletes blobs by key.
type Storage interface {
	Upload(ctx context.Context, key string, content io.Reader) error
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
}

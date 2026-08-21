// Package secrets defines the common Secrets port. All option subpackages
// implement this interface. env is the always-available default; cloud
// secret managers (gcp, aws, vault) are additive.
package secrets

import "context"

// Secrets retrieves secret values by key.
type Secrets interface {
	Get(ctx context.Context, key string) (string, error)
}

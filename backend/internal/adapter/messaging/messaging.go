// Package messaging defines the common Publisher port. All option
// subpackages implement this interface.
package messaging

import "context"

// Publisher emits messages to a broker. Implementations are in
// internal/adapter/messaging/<option>/.
type Publisher interface {
	Publish(ctx context.Context, topic string, payload []byte) error
	Close() error
}

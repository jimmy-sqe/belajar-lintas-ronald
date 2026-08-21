// Package notification defines the common Notification port. All option
// subpackages implement this interface.
package notification

import "context"

// Notification sends transactional messages.
type Notification interface {
	Send(ctx context.Context, to, subject, body string) error
}

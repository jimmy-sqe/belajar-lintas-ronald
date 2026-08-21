// Package noop is the no-op observability adapter. It does nothing —
// useful for environments without observability infrastructure.
package noop

import "context"

// Config has no fields.
type Config struct{}

// Setup is a no-op that returns a no-op shutdown function.
func Setup(_ context.Context, _ Config) (func(context.Context) error, error) {
	return func(_ context.Context) error { return nil }, nil
}

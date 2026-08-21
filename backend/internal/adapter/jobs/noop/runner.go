// Package noop is the no-op jobs adapter. Schedule/Start/Stop are silent.
package noop

import "context"

// Runner is a no-op Runner.
type Runner struct{}

// New returns a noop Runner.
func New() *Runner { return &Runner{} }

// Schedule silently accepts any spec.
func (r *Runner) Schedule(_ string, _ func(ctx context.Context)) error { return nil }

// Start is a no-op.
func (r *Runner) Start() {}

// Stop is a no-op.
func (r *Runner) Stop(_ context.Context) error { return nil }

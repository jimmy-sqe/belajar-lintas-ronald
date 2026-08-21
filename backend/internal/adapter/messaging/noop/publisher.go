// Package noop is the no-op messaging adapter. Publish is a silent success.
package noop

import (
	"context"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/messaging"
)

// Publisher is the no-op implementation.
type Publisher struct{}

// New returns a no-op Publisher.
func New() messaging.Publisher { return &Publisher{} }

// Publish silently discards.
func (p *Publisher) Publish(_ context.Context, _ string, _ []byte) error { return nil }

// Close is a no-op.
func (p *Publisher) Close() error { return nil }

// Package noop is the no-op notification adapter. Send is a silent success.
package noop

import (
	"context"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/notification"
)

// Client is the no-op implementation.
type Client struct{}

// New returns a noop Client.
func New() notification.Notification { return &Client{} }

// Send silently discards.
func (c *Client) Send(_ context.Context, _, _, _ string) error { return nil }

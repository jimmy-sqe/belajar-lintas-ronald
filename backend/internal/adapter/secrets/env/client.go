// Package env is the OS-environment-variables secrets adapter.
package env

import (
	"context"
	"fmt"
	"os"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/secrets"
)

// Client reads secrets from the process environment.
type Client struct{}

// New returns an env Client.
func New() secrets.Secrets { return &Client{} }

// Get returns the value of the named env var, or an error if unset.
func (c *Client) Get(_ context.Context, key string) (string, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("secrets/env: %s is not set", key)
	}
	return v, nil
}

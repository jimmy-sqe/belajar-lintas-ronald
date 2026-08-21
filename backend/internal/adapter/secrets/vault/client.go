// Package vault is the HashiCorp Vault secrets adapter.
package vault

import (
	"context"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// Config holds Vault settings.
type Config struct {
	Address string `mapstructure:"VAULT_ADDR"`
	Token   string `mapstructure:"VAULT_TOKEN"`
	Path    string `mapstructure:"VAULT_KV_PATH"` // KV mount path, e.g. "secret/data"
}

// Client wraps a Vault API client.
type Client struct {
	client *vaultapi.Client
	path   string
}

// New constructs a Client.
func New(_ context.Context, cfg Config) (*Client, error) {
	if cfg.Address == "" || cfg.Token == "" {
		return nil, fmt.Errorf("secrets/vault: address and token are required")
	}
	apiCfg := vaultapi.DefaultConfig()
	apiCfg.Address = cfg.Address
	c, err := vaultapi.NewClient(apiCfg)
	if err != nil {
		return nil, fmt.Errorf("secrets/vault: client: %w", err)
	}
	c.SetToken(cfg.Token)
	path := cfg.Path
	if path == "" {
		path = "secret/data"
	}
	return &Client{client: c, path: path}, nil
}

// Get reads the named secret from the KV mount; returns the "value" field
// of the secret data (KV v2 wraps data inside a "data" map).
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	sec, err := c.client.Logical().ReadWithContext(ctx, fmt.Sprintf("%s/%s", c.path, key))
	if err != nil {
		return "", fmt.Errorf("secrets/vault: read: %w", err)
	}
	if sec == nil || sec.Data == nil {
		return "", fmt.Errorf("secrets/vault: %s not found", key)
	}
	// KV v2: actual data is under sec.Data["data"]
	if data, ok := sec.Data["data"].(map[string]any); ok {
		if v, ok := data["value"].(string); ok {
			return v, nil
		}
		return "", fmt.Errorf("secrets/vault: %s has no string 'value' field", key)
	}
	// KV v1 fallback
	if v, ok := sec.Data["value"].(string); ok {
		return v, nil
	}
	return "", fmt.Errorf("secrets/vault: %s payload format unknown", key)
}

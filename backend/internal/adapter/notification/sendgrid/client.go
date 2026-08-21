// Package sendgrid is the SendGrid notification adapter.
package sendgrid

import (
	"context"
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Config holds SendGrid settings.
type Config struct {
	APIKey string `mapstructure:"SENDGRID_API_KEY"`
	From   string `mapstructure:"SENDGRID_FROM"`
}

// Client wraps the SendGrid REST client.
type Client struct {
	client *sendgrid.Client
	from   string
}

// New constructs a Client.
func New(cfg Config) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("sendgrid: API key is required")
	}
	if cfg.From == "" {
		return nil, fmt.Errorf("sendgrid: from is required")
	}
	return &Client{client: sendgrid.NewSendClient(cfg.APIKey), from: cfg.From}, nil
}

// Send dispatches a plain-text email via SendGrid.
func (c *Client) Send(ctx context.Context, to, subject, body string) error {
	from := mail.NewEmail("", c.from)
	toAddr := mail.NewEmail("", to)
	msg := mail.NewSingleEmail(from, subject, toAddr, body, body)
	resp, err := c.client.SendWithContext(ctx, msg)
	if err != nil {
		return fmt.Errorf("sendgrid: send: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("sendgrid: status %d: %s", resp.StatusCode, resp.Body)
	}
	return nil
}

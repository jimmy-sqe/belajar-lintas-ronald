// Package smtp is the SMTP notification adapter using net/smtp.
package smtp

import (
	"context"
	"fmt"
	"net/smtp"
)

// Config holds SMTP settings.
type Config struct {
	Host     string `mapstructure:"SMTP_HOST"`
	Port     int    `mapstructure:"SMTP_PORT"`
	Username string `mapstructure:"SMTP_USERNAME"`
	Password string `mapstructure:"SMTP_PASSWORD"`
	From     string `mapstructure:"SMTP_FROM"`
}

// Client is the SMTP implementation.
type Client struct {
	cfg Config
}

// New returns a Client. If From is empty, errors at Send time.
func New(cfg Config) (*Client, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("smtp: host is required")
	}
	if cfg.Port == 0 {
		cfg.Port = 587
	}
	return &Client{cfg: cfg}, nil
}

// Send dispatches a plain-text email via SMTP using PLAIN auth.
func (c *Client) Send(_ context.Context, to, subject, body string) error {
	if c.cfg.From == "" {
		return fmt.Errorf("smtp: From is required")
	}
	addr := fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port)
	auth := smtp.PlainAuth("", c.cfg.Username, c.cfg.Password, c.cfg.Host)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		c.cfg.From, to, subject, body)
	if err := smtp.SendMail(addr, auth, c.cfg.From, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("smtp: send: %w", err)
	}
	return nil
}

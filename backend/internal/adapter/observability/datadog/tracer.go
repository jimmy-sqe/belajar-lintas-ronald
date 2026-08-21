// Package datadog is the Datadog observability adapter. It uses dd-trace-go
// to send traces to a Datadog agent.
package datadog

import (
	"context"
	"fmt"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Config holds Datadog tracer settings.
type Config struct {
	AgentHost   string `mapstructure:"DD_AGENT_HOST"`
	AgentPort   string `mapstructure:"DD_TRACE_AGENT_PORT"`
	ServiceName string `mapstructure:"DD_SERVICE"`
	Environment string `mapstructure:"DD_ENV"`
	Enabled     bool   `mapstructure:"DD_TRACE_ENABLED"`
}

// Setup starts the Datadog tracer. Returned shutdown func should be called
// on app exit.
func Setup(_ context.Context, cfg Config) (func(context.Context) error, error) {
	if !cfg.Enabled {
		return func(_ context.Context) error { return nil }, nil
	}
	if cfg.ServiceName == "" {
		cfg.ServiceName = "backend-belajar-lintas-ronald"
	}
	opts := []tracer.StartOption{
		tracer.WithService(cfg.ServiceName),
	}
	if cfg.Environment != "" {
		opts = append(opts, tracer.WithEnv(cfg.Environment))
	}
	if cfg.AgentHost != "" {
		port := cfg.AgentPort
		if port == "" {
			port = "8126"
		}
		opts = append(opts, tracer.WithAgentAddr(fmt.Sprintf("%s:%s", cfg.AgentHost, port)))
	}
	tracer.Start(opts...)
	return func(_ context.Context) error {
		tracer.Stop()
		return nil
	}, nil
}

package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/config"
)

func TestNewObservability_Noop_FRObs(t *testing.T) {
	cfg := &config.Config{ObservabilityBackend: "noop"}

	obs, err := newObservability(context.Background(), cfg)

	require.NoError(t, err)
	require.NotNil(t, obs)
	require.NotNil(t, obs.middleware)
	require.NotNil(t, obs.shutdown)
	require.NoError(t, obs.shutdown(context.Background()))
}

func TestNewObservability_Unknown_FRObs(t *testing.T) {
	cfg := &config.Config{ObservabilityBackend: "bogus"}

	_, err := newObservability(context.Background(), cfg)

	require.Error(t, err)
}

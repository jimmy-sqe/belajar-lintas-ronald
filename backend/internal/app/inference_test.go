package app

import (
	"testing"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/config"

	"github.com/stretchr/testify/require"
)

func TestNewInference_Noop(t *testing.T) {
	cfg := &config.Config{InferenceBackend: "noop"}
	eng, err := newInference(cfg)
	require.NoError(t, err)
	require.NotNil(t, eng)
}

func TestNewInference_Unknown(t *testing.T) {
	cfg := &config.Config{InferenceBackend: "bogus"}
	_, err := newInference(cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "INFERENCE_BACKEND")
}

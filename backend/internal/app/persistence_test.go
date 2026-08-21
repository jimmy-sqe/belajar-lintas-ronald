package app

import (
	"context"
	"testing"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/config"

	"github.com/stretchr/testify/require"
)

func TestNewTodoRepo_UnknownBackend(t *testing.T) {
	cfg := &config.Config{PersistenceBackend: "noop"}
	_, _, err := newTodoRepo(context.Background(), cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "PERSISTENCE_BACKEND")
}

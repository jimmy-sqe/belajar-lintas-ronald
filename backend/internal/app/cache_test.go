package app

import (
	"context"
	"testing"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/config"

	"github.com/stretchr/testify/require"
)

func TestNewTodoCache_Inmemory(t *testing.T) {
	cfg := &config.Config{CacheBackend: "inmemory"}
	c, closeFn, err := newTodoCache(context.Background(), cfg)
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotNil(t, closeFn)
}

func TestNewTodoCache_Noop(t *testing.T) {
	cfg := &config.Config{CacheBackend: "noop"}
	c, _, err := newTodoCache(context.Background(), cfg)
	require.NoError(t, err)
	require.NotNil(t, c)
}

func TestNewTodoCache_Unknown(t *testing.T) {
	cfg := &config.Config{CacheBackend: "bogus"}
	_, _, err := newTodoCache(context.Background(), cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "CACHE_BACKEND")
}

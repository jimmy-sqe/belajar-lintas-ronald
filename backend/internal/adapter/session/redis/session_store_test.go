//go:build integration
// +build integration

package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"

	authdomain "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/auth"
	rsessions "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/session/redis"
)

func TestSessionStore_PutGetDelete(t *testing.T) {
	ctx := context.Background()
	container, err := tcredis.Run(ctx, "redis:7-alpine")
	require.NoError(t, err)
	t.Cleanup(func() { _ = testcontainers.TerminateContainer(container) })

	conn, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	store, err := rsessions.NewSessionStore(conn)
	require.NoError(t, err)

	userID := uuid.New()
	hash := "deadbeef"
	expiry := time.Now().UTC().Add(time.Hour)

	require.NoError(t, store.Put(ctx, authdomain.RefreshSession{
		TokenHash: hash,
		UserID:    userID,
		ExpiresAt: expiry,
	}))

	got, err := store.Get(ctx, hash)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, userID, got.UserID)
	assert.WithinDuration(t, expiry, got.ExpiresAt, time.Second)

	require.NoError(t, store.Delete(ctx, hash))

	gone, err := store.Get(ctx, hash)
	require.NoError(t, err)
	assert.Nil(t, gone)
}

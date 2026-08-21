//go:build integration

package cmd_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcmongo "github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/cmd"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/persistence/mongodb"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/config"
)

func setupMongoForSeed(t *testing.T) (*config.Config, func()) {
	t.Helper()
	ctx := context.Background()
	c, err := tcmongo.RunContainer(ctx, testcontainers.WithImage("mongo:7"))
	require.NoError(t, err)
	uri, err := c.ConnectionString(ctx)
	require.NoError(t, err)

	cfg := &config.Config{
		Mongo: mongodb.Config{URI: uri, DB: "test"},
	}
	cleanup := func() { _ = c.Terminate(ctx) }
	return cfg, cleanup
}

func TestSeed_Mongo_InsertsTwoSampleDocs(t *testing.T) {
	chdirToGolangRoot(t)
	cfg, done := setupMongoForSeed(t)
	defer done()

	ctx := context.Background()
	require.NoError(t, cmd.RunSeeds(ctx, cfg))

	client, err := mongodb.Connect(ctx, cfg.Mongo)
	require.NoError(t, err)
	defer func() { _ = client.Disconnect(context.Background()) }()
	count, err := client.Database("test").Collection("items").CountDocuments(ctx, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(2), count)
}

func TestSeed_Mongo_Idempotent(t *testing.T) {
	chdirToGolangRoot(t)
	cfg, done := setupMongoForSeed(t)
	defer done()

	ctx := context.Background()
	require.NoError(t, cmd.RunSeeds(ctx, cfg))
	require.NoError(t, cmd.RunSeeds(ctx, cfg))

	client, err := mongodb.Connect(ctx, cfg.Mongo)
	require.NoError(t, err)
	defer func() { _ = client.Disconnect(context.Background()) }()
	count, err := client.Database("test").Collection("items").CountDocuments(ctx, bson.M{})
	require.NoError(t, err)
	require.Equal(t, int64(2), count)
}

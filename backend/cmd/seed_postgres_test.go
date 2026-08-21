//go:build integration

package cmd_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migpg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/cmd"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/persistence/postgres"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/config"
)

// chdirToGolangRoot moves CWD to the golang/ subproject root so that the seed
// executor's relative paths (db/seeds/postgres) resolve.
func chdirToGolangRoot(t *testing.T) {
	t.Helper()
	abs, err := filepath.Abs("..") // cmd/ -> golang/
	require.NoError(t, err)
	t.Chdir(abs)
}

func setupPostgresForSeed(t *testing.T) (*config.Config, func()) {
	t.Helper()
	ctx := context.Background()

	pg, err := tcpg.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		tcpg.WithDatabase("test"),
		tcpg.WithUsername("test"),
		tcpg.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(60*time.Second)),
	)
	require.NoError(t, err)

	host, _ := pg.Host(ctx)
	port, _ := pg.MappedPort(ctx, "5432")

	cfg := &config.Config{
		Postgres: postgres.Config{
			Host: host, Port: int(port.Num()), User: "test",
			Password: "test", DB: "test", SSLMode: "disable",
		},
	}

	// Apply schema migrations so the items table exists.
	db, err := postgres.Connect(ctx, cfg.Postgres)
	require.NoError(t, err)
	driver, err := migpg.WithInstance(db.DB, &migpg.Config{})
	require.NoError(t, err)
	migPath, _ := filepath.Abs("db/migrations/postgres")
	m, err := migrate.NewWithDatabaseInstance("file://"+migPath, "postgres", driver)
	require.NoError(t, err)
	require.NoError(t, m.Up())

	cleanup := func() {
		_ = db.Close()
		_ = pg.Terminate(ctx)
	}
	return cfg, cleanup
}

func TestSeed_Postgres_InsertsTwoSampleRows(t *testing.T) {
	chdirToGolangRoot(t)
	cfg, done := setupPostgresForSeed(t)
	defer done()

	ctx := context.Background()
	require.NoError(t, cmd.RunSeeds(ctx, cfg))

	db, err := postgres.Connect(ctx, cfg.Postgres)
	require.NoError(t, err)
	defer func() { _ = db.Close() }()
	var count int
	require.NoError(t, db.GetContext(ctx, &count, "SELECT COUNT(*) FROM items"))
	require.Equal(t, 2, count)
}

func TestSeed_Postgres_Idempotent(t *testing.T) {
	chdirToGolangRoot(t)
	cfg, done := setupPostgresForSeed(t)
	defer done()

	ctx := context.Background()
	require.NoError(t, cmd.RunSeeds(ctx, cfg))
	require.NoError(t, cmd.RunSeeds(ctx, cfg)) // second run must not error

	db, err := postgres.Connect(ctx, cfg.Postgres)
	require.NoError(t, err)
	defer func() { _ = db.Close() }()
	var count int
	require.NoError(t, db.GetContext(ctx, &count, "SELECT COUNT(*) FROM items"))
	require.Equal(t, 2, count, "idempotency: re-running seed must not duplicate rows")
}

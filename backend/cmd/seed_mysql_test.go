//go:build integration

package cmd_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	migmysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/cmd"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/persistence/mysql"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/config"
)

func setupMySQLForSeed(t *testing.T) (*config.Config, func()) {
	t.Helper()
	ctx := context.Background()

	c, err := tcmysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:8.0"),
		tcmysql.WithDatabase("test"),
		tcmysql.WithUsername("test"),
		tcmysql.WithPassword("test"),
	)
	require.NoError(t, err)

	host, _ := c.Host(ctx)
	port, _ := c.MappedPort(ctx, "3306")

	cfg := &config.Config{
		MySQL: mysql.Config{
			Host: host, Port: int(port.Num()), User: "test",
			Password: "test", DB: "test",
		},
	}

	db, err := mysql.Connect(ctx, cfg.MySQL)
	require.NoError(t, err)
	driver, err := migmysql.WithInstance(db.DB, &migmysql.Config{})
	require.NoError(t, err)
	migPath, _ := filepath.Abs("db/migrations/mysql")
	m, err := migrate.NewWithDatabaseInstance("file://"+migPath, "mysql", driver)
	require.NoError(t, err)
	require.NoError(t, m.Up())

	cleanup := func() {
		_ = db.Close()
		_ = c.Terminate(ctx)
	}
	return cfg, cleanup
}

func TestSeed_MySQL_InsertsTwoSampleRows(t *testing.T) {
	chdirToGolangRoot(t)
	cfg, done := setupMySQLForSeed(t)
	defer done()

	ctx := context.Background()
	require.NoError(t, cmd.RunSeeds(ctx, cfg))

	db, err := mysql.Connect(ctx, cfg.MySQL)
	require.NoError(t, err)
	defer func() { _ = db.Close() }()
	var count int
	require.NoError(t, db.GetContext(ctx, &count, "SELECT COUNT(*) FROM items"))
	require.Equal(t, 2, count)
}

func TestSeed_MySQL_Idempotent(t *testing.T) {
	chdirToGolangRoot(t)
	cfg, done := setupMySQLForSeed(t)
	defer done()

	ctx := context.Background()
	require.NoError(t, cmd.RunSeeds(ctx, cfg))
	require.NoError(t, cmd.RunSeeds(ctx, cfg))

	db, err := mysql.Connect(ctx, cfg.MySQL)
	require.NoError(t, err)
	defer func() { _ = db.Close() }()
	var count int
	require.NoError(t, db.GetContext(ctx, &count, "SELECT COUNT(*) FROM items"))
	require.Equal(t, 2, count)
}

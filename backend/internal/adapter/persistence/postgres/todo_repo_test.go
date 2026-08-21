//go:build integration

package postgres_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migpg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/persistence/postgres"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/todo"
)

func setupDB(t *testing.T) (*sqlx.DB, func()) {
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

	db, err := postgres.Connect(ctx, postgres.Config{
		Host:     host,
		Port:     int(port.Num()),
		User:     "test",
		Password: "test",
		DB:       "test",
		SSLMode:  "disable",
	})
	require.NoError(t, err)

	// Run migrations
	driver, err := migpg.WithInstance(db.DB, &migpg.Config{})
	require.NoError(t, err)
	migrationsPath, err := filepath.Abs("../../../../db/migrations/postgres")
	require.NoError(t, err)
	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsPath, "postgres", driver)
	require.NoError(t, err)
	require.NoError(t, m.Up())

	cleanup := func() {
		_ = db.Close()
		_ = pg.Terminate(ctx)
	}
	return db, cleanup
}

func TestTodoRepository_CreateAndFindByID(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()

	repo := postgres.NewTodoRepository(db)
	ctx := context.Background()

	ownerID := uuid.New()
	now := time.Now().UTC().Truncate(time.Microsecond)
	in := &todo.Todo{
		ID:         uuid.New(),
		Title:      "Test Todo",
		Description: "Test description",
		CreatedAt:  now,
		CreatedBy:  ownerID,
		ModifiedAt: now,
		ModifiedBy: ownerID,
	}
	require.NoError(t, repo.Create(ctx, in))

	got, err := repo.FindByID(ctx, in.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, in.Title, got.Title)
	require.Equal(t, in.Description, got.Description)
	require.Equal(t, in.CreatedBy, got.CreatedBy)
}

func TestTodoRepository_FindByID_NotFound(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()

	repo := postgres.NewTodoRepository(db)
	got, err := repo.FindByID(context.Background(), uuid.New())
	require.ErrorIs(t, err, todo.ErrNotFound)
	require.Nil(t, got)
}

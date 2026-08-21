//go:build integration

package mysql_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migmysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/persistence/mysql"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/todo"
)

func setupDB(t *testing.T) (*sqlx.DB, func()) {
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

	db, err := mysql.Connect(ctx, mysql.Config{
		Host: host, Port: int(port.Num()), User: "test",
		Password: "test", DB: "test",
	})
	require.NoError(t, err)

	driver, err := migmysql.WithInstance(db.DB, &migmysql.Config{})
	require.NoError(t, err)
	migrationsPath, err := filepath.Abs("../../../../db/migrations/mysql")
	require.NoError(t, err)
	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsPath, "mysql", driver)
	require.NoError(t, err)
	require.NoError(t, m.Up())

	cleanup := func() {
		_ = db.Close()
		_ = c.Terminate(ctx)
	}
	return db, cleanup
}

func TestTodoRepository_CreateAndFindByID(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()

	repo := mysql.NewTodoRepository(db)
	ctx := context.Background()

	ownerID := uuid.New()
	now := time.Now().UTC().Truncate(time.Microsecond)
	in := &todo.Todo{
		ID:          uuid.New(),
		Title:       "Test Todo",
		Description: "Test description",
		CreatedAt:   now,
		CreatedBy:   ownerID,
		ModifiedAt:  now,
		ModifiedBy:  ownerID,
	}
	require.NoError(t, repo.Create(ctx, in))

	got, err := repo.FindByID(ctx, in.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, in.Title, got.Title)
}

func TestTodoRepository_FindByID_NotFound(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()

	repo := mysql.NewTodoRepository(db)
	got, err := repo.FindByID(context.Background(), uuid.New())
	require.ErrorIs(t, err, todo.ErrNotFound)
	require.Nil(t, got)
}

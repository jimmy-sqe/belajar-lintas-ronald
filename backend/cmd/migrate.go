package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migpg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	// boilerplate:axis=persistence option=postgres START
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/persistence/postgres"
	// boilerplate:axis=persistence option=postgres END

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/config"
)

func migrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db:migrate",
		Short: "Run database migrations up",
		RunE: func(c *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			return runMigrations(c.Context(), cfg)
		},
	}
	return cmd
}

func runMigrations(ctx context.Context, cfg *config.Config) error {
	// boilerplate:axis=persistence option=postgres START
	db, err := postgres.Connect(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer db.Close() //nolint:errcheck // best-effort close on shutdown
	if err := migrateOne(db, "db/migrations/postgres", "postgres"); err != nil {
		return err
	}
	// boilerplate:axis=persistence option=postgres END
	return nil
}

func migrateOne(db *sqlx.DB, path, dbName string) error {
	driver, err := migpg.WithInstance(db.DB, &migpg.Config{})
	if err != nil {
		return fmt.Errorf("driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+path, dbName, driver)
	if err != nil {
		return fmt.Errorf("new migrate: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}

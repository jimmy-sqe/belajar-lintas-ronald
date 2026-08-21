package app

import (
	"context"
	"fmt"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/config"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/todo"

	// boilerplate:axis=persistence option=postgres START
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/persistence/postgres"
	// boilerplate:axis=persistence option=postgres END
	// boilerplate:axis=persistence option=mysql START
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/persistence/mysql"
	// boilerplate:axis=persistence option=mysql END
	// boilerplate:axis=persistence option=mongodb START
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/persistence/mongodb"
	// boilerplate:axis=persistence option=mongodb END
)

// newTodoRepo builds the todo Repository for the selected persistence backend.
// Mirrors newStorage: per-option marker-wrapped, unknown value fails fast.
// There is intentionally no noop arm — persistence=noop has no Repository consumer.
func newTodoRepo(ctx context.Context, cfg *config.Config) (todo.Repository, func(context.Context) error, error) {
	var valid []string
	// boilerplate:axis=persistence option=postgres START
	valid = append(valid, "postgres")
	// boilerplate:axis=persistence option=postgres END
	// boilerplate:axis=persistence option=mysql START
	valid = append(valid, "mysql")
	// boilerplate:axis=persistence option=mysql END
	// boilerplate:axis=persistence option=mongodb START
	valid = append(valid, "mongodb")
	// boilerplate:axis=persistence option=mongodb END

	switch cfg.PersistenceBackend {
	// boilerplate:axis=persistence option=postgres START
	case "postgres":
		db, err := postgres.Connect(ctx, cfg.Postgres)
		if err != nil {
			return nil, nil, fmt.Errorf("app: postgres: %w", err)
		}
		return postgres.NewTodoRepository(db), func(_ context.Context) error { return db.Close() }, nil
	// boilerplate:axis=persistence option=postgres END
	// boilerplate:axis=persistence option=mysql START
	case "mysql":
		db, err := mysql.Connect(ctx, cfg.MySQL)
		if err != nil {
			return nil, nil, fmt.Errorf("app: mysql: %w", err)
		}
		return mysql.NewTodoRepository(db), func(_ context.Context) error { return db.Close() }, nil
	// boilerplate:axis=persistence option=mysql END
	// boilerplate:axis=persistence option=mongodb START
	case "mongodb":
		client, err := mongodb.Connect(ctx, cfg.Mongo)
		if err != nil {
			return nil, nil, fmt.Errorf("app: mongodb: %w", err)
		}
		repo := mongodb.NewTodoRepository(client.Database(cfg.Mongo.DB))
		return repo, func(c context.Context) error { return client.Disconnect(c) }, nil
	// boilerplate:axis=persistence option=mongodb END
	}
	return nil, nil, fmt.Errorf("persistence: PERSISTENCE_BACKEND must be one of %v; got %q", valid, cfg.PersistenceBackend)
}

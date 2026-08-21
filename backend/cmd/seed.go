package cmd

import (
	"context"
	// boilerplate:axis=persistence option=mongodb START
	"encoding/json"
	// boilerplate:axis=persistence option=mongodb END
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	// boilerplate:axis=persistence option=mongodb START
	"time"
	// boilerplate:axis=persistence option=mongodb END

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	// boilerplate:axis=persistence option=mongodb START
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// boilerplate:axis=persistence option=mongodb END

	// boilerplate:axis=persistence option=postgres START
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/persistence/postgres"
	// boilerplate:axis=persistence option=postgres END

	// boilerplate:axis=persistence option=mysql START
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/persistence/mysql"
	// boilerplate:axis=persistence option=mysql END

	// boilerplate:axis=persistence option=mongodb START
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/persistence/mongodb"
	// boilerplate:axis=persistence option=mongodb END

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/config"
)

const (
	// boilerplate:axis=persistence option=postgres START
	seedDirPostgres = "db/seeds/postgres"
	// boilerplate:axis=persistence option=postgres END
	// boilerplate:axis=persistence option=mysql START
	seedDirMySQL = "db/seeds/mysql"
	// boilerplate:axis=persistence option=mysql END
	// boilerplate:axis=persistence option=mongodb START
	seedDirMongo = "db/seeds/mongodb"
	// boilerplate:axis=persistence option=mongodb END
)

func seedCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "db:seed",
		Short: "Apply seed data for the selected persistence dialect",
		RunE: func(c *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			return RunSeeds(c.Context(), cfg)
		},
	}
}

// RunSeeds applies seed data to whichever persistence dialects are configured.
// Each dialect's block is wrapped in axis markers so the SDLC plugin prunes
// unselected dialects.
func RunSeeds(ctx context.Context, cfg *config.Config) error {
	// boilerplate:axis=persistence option=postgres START
	if cfg.Postgres.Host != "" {
		db, err := postgres.Connect(ctx, cfg.Postgres)
		if err != nil {
			return fmt.Errorf("seed/postgres: connect: %w", err)
		}
		defer func() { _ = db.Close() }()
		if err := execSQLDir(ctx, db, seedDirPostgres); err != nil {
			return err
		}
	}
	// boilerplate:axis=persistence option=postgres END

	// boilerplate:axis=persistence option=mysql START
	if cfg.MySQL.Host != "" {
		db, err := mysql.Connect(ctx, cfg.MySQL)
		if err != nil {
			return fmt.Errorf("seed/mysql: connect: %w", err)
		}
		defer func() { _ = db.Close() }()
		if err := execSQLDir(ctx, db, seedDirMySQL); err != nil {
			return err
		}
	}
	// boilerplate:axis=persistence option=mysql END

	// boilerplate:axis=persistence option=mongodb START
	if cfg.Mongo.URI != "" {
		client, err := mongodb.Connect(ctx, cfg.Mongo)
		if err != nil {
			return fmt.Errorf("seed/mongodb: connect: %w", err)
		}
		defer func() { _ = client.Disconnect(context.Background()) }()
		if err := execJSONDir(ctx, client.Database(cfg.Mongo.DB), seedDirMongo); err != nil {
			return err
		}
	}
	// boilerplate:axis=persistence option=mongodb END

	return nil
}

// execSQLDir runs every *.sql file in dir in lexicographic order against db.
// Each file is treated as one logical statement.
func execSQLDir(ctx context.Context, db *sqlx.DB, dir string) error {
	files, err := listFiles(dir, ".sql")
	if err != nil {
		return err
	}
	for _, f := range files {
		content, err := os.ReadFile(f) //nolint:gosec // path constrained to seed dir constant
		if err != nil {
			return fmt.Errorf("seed: read %s: %w", f, err)
		}
		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("seed: exec %s: %w", f, err)
		}
	}
	return nil
}

// boilerplate:axis=persistence option=mongodb START
// execJSONDir runs every *.json file in dir. Each file's name (minus NNN_
// prefix) is the target collection. Each document is upserted via $setOnInsert.
func execJSONDir(ctx context.Context, db *mongo.Database, dir string) error {
	files, err := listFiles(dir, ".json")
	if err != nil {
		return err
	}
	for _, f := range files {
		content, err := os.ReadFile(f) //nolint:gosec // path constrained to seed dir constant
		if err != nil {
			return fmt.Errorf("seed: read %s: %w", f, err)
		}
		var docs []bson.M
		if err := json.Unmarshal(content, &docs); err != nil {
			return fmt.Errorf("seed: parse %s: %w", f, err)
		}
		coll := db.Collection(collectionFromFilename(f))
		for i, doc := range docs {
			convertDateStrings(doc)
			id, ok := doc["_id"]
			if !ok {
				return fmt.Errorf("seed: %s doc[%d] missing _id", f, i)
			}
			if _, err := coll.UpdateOne(ctx,
				bson.M{"_id": id},
				bson.M{"$setOnInsert": doc},
				options.Update().SetUpsert(true),
			); err != nil {
				return fmt.Errorf("seed: upsert %s doc[%d]: %w", f, i, err)
			}
		}
	}
	return nil
}

// boilerplate:axis=persistence option=mongodb END

// listFiles returns sorted absolute paths of files in dir with the given extension.
// Returns nil if dir does not exist (treated as "nothing to seed").
func listFiles(dir, ext string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("seed: list %s: %w", dir, err)
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ext) {
			continue
		}
		out = append(out, filepath.Join(dir, e.Name()))
	}
	sort.Strings(out)
	return out, nil
}

// boilerplate:axis=persistence option=mongodb START

// collectionFromFilename derives the target collection from a filename of the
// form NNN_<collection>.json. Falls back to the trimmed basename if there is
// no numeric prefix.
func collectionFromFilename(path string) string {
	name := strings.TrimSuffix(filepath.Base(path), ".json")
	if idx := strings.Index(name, "_"); idx > 0 && idx <= 4 {
		prefix := name[:idx]
		allDigits := true
		for _, r := range prefix {
			if r < '0' || r > '9' {
				allDigits = false
				break
			}
		}
		if allDigits {
			return name[idx+1:]
		}
	}
	return name
}

// convertDateStrings replaces RFC3339 strings in known datetime fields with
// time.Time values so they land in MongoDB as BSON Date, not strings.
func convertDateStrings(doc bson.M) {
	for _, field := range []string{"created_at", "updated_at", "deleted_at"} {
		v, ok := doc[field]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		t, err := time.Parse(time.RFC3339, s)
		if err == nil {
			doc[field] = t
		}
	}
}

// boilerplate:axis=persistence option=mongodb END

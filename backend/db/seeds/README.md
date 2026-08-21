# Database Seeds

This directory holds initial / master data loaded by `make seed-<dialect>` (which runs `go run . db:seed`).

## Rules

1. **Every seed file MUST be idempotent.** Re-running the seed command must not duplicate rows or fail.
2. **One logical statement per `.sql` file.** Authors needing multiple statements split them into multiple `NNN_*.sql` files.
3. **Naming:** `NNN_<name>.{sql,json}` where `NNN` is a zero-padded sequence (001, 002, ...). Order is alphabetical by filename.
4. **MongoDB filename = target collection.** `001_todos.json` loads into the `todos` collection. To seed multiple collections, use multiple files (`001_todos.json`, `002_users.json`).

## Idempotency patterns

| Dialect  | Pattern |
|----------|---------|
| Postgres | `INSERT ... ON CONFLICT (id) DO NOTHING` |
| MySQL    | `INSERT IGNORE INTO ...` |
| MongoDB  | Per-document `UpdateOne` with `$setOnInsert` + `upsert: true` (handled by the Go executor) |

The Go executor itself does not deduplicate. Authors are responsible for writing idempotent statements.

## How to add a new seed

1. Pick the next sequence number for the target dialect (`ls db/seeds/<dialect>/ | tail -1`).
2. Create the file using the patterns above.
3. Run `make seed-<dialect>` to verify it applies cleanly.

## How to reset seed data

Out of scope for the boilerplate. Operators reset manually via `TRUNCATE` / `db.collection.drop()` and then re-run migrations + seeds.

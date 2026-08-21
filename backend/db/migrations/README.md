# Database Migrations

This directory holds schema migrations applied by `make migrate-<dialect>` (which runs `go run . db:migrate`).

## Policy: forward-only

The boilerplate adopts a **forward-only migration policy**. There are no `.down.sql` files.

### Why no down?

- In production, accidentally running `migrate down` is a permanent data-loss hazard (e.g., `DROP TABLE` reverses a `CREATE TABLE` and erases data).
- Modern operational practice favours fix-forward over rollback for schema changes: if a migration introduces a bug, write a new migration that inverts it. The history remains a linear sequence of additive changes.

### Local development rollback

If you need to undo a schema change while iterating locally:

1. Drop the schema (or the affected table).
2. Re-run `make migrate-<dialect>` from a clean state.

Local databases are ephemeral; re-seeding is fast.

## Naming

- Postgres / MySQL: `NNN_<description>.up.sql` (only `.up.sql` — no `.down.sql`).
- MongoDB: schema-on-write; no migrations here today. Add Go-based migrators if/when needed.

## How to add a migration

1. Find the highest existing prefix: `ls db/migrations/postgres/ | tail -1`.
2. Pick the next number.
3. Write the change as a single `.up.sql` file.
4. Verify locally with `make migrate-postgres` against a fresh DB.
5. Commit; CI runs the same migration against a containerized DB during integration tests.

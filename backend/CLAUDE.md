# Go Service — Working Rules

This file applies when working in this Go service (boilerplate or
scaffolded). Loaded automatically by Claude Code.

## Stack at a glance

Go 1.25 · Echo v4 · sqlx + Squirrel · Viper · zerolog · testify +
testcontainers. Architecture: Hexagonal (Ports & Adapters).

→ Full architecture: `docs/ARCHITECTURE.md`
→ Make targets: `README.md#tooling`
→ Manifest: `.boilerplate.yaml`

## Rule 1 — Hexagonal layering

Strict dependency direction enforced by `golangci-lint depguard`:

- Domain (`internal/domain/<feature>/`) declares ports as interfaces.
  Pure business logic. May import only `pkg/*`. May NOT import
  `internal/adapter/*`, `internal/http/*`, `internal/app/*`.
- Adapters (`internal/adapter/<axis>/<option>/`) implement domain ports.
  Each adapter may NOT import another adapter.
- Inbound HTTP (`internal/http/`) translates HTTP ↔ domain via service
  interfaces. May import `internal/adapter/auth/*` (middleware bridges
  the HTTP boundary) but not other adapter types.
- Composition root (`internal/app/container.go`) is the only file that
  imports many packages — it wires concrete adapters into domain
  services.
- Cross-domain access goes through the other domain's `Service`
  interface, never its `Repository`. The service carries business
  rules (authz, cache, audit, events) you must not bypass.

→ Detail + diagrams: `docs/ARCHITECTURE.md`

## Rule 2 — Response envelope (lintas-standard)

Every HTTP response (except `/healthz`) uses helpers from
`internal/http/response/`:

- `response.OK(c, http.StatusOK, "message", data)` — single-resource success
- `response.Paged(c, "message", data, pagination)` — list with pagination
- `response.Err(c, err)` — error (accepts any `error`; non-custom errors render as 500)

Shape:

```json
// Success
{"success": true, "code": 200, "message": "...", "data": {...}, "timestamp": "..."}
// Paginated (adds "pagination")
{"success": true, "code": 200, "message": "...", "data": [...], "pagination": {...}, "timestamp": "..."}
// Error
{"success": false, "code": 40400, "message": "...", "timestamp": "...", "metadata": {...}}
```

- `code` is a 5-digit integer (`HTTP_PREFIX + 2-digit subcode`).
  Subcode `00` = generic; `01-99` = service-specific.
- DELETE endpoints return `200 OK` + `data: {}` — NOT 204.
- Handlers MUST go through these helpers. Manual `c.JSON(...)` of
  envelope-shaped maps defeats the single-source-of-truth guarantee.

→ Full error catalog + how to extend: `docs/ERROR_CODES.md`
→ Envelope contract: `docs/ARCHITECTURE.md#response-envelope`

## Rule 3 — FR annotation pattern

Every handler, service method, or repository method that implements a
functional requirement carries a `// FR-<ID>: <one-line description>`
comment above its declaration. Lintas QA and code-review skills parse
this annotation to compute coverage.

```go
// FR-001: List todos for the authenticated user, paginated.
func (h *TodoHandler) List(c echo.Context) error { ... }
```

→ Annotation parser: lintas plugin `sdlc-qa:qa-validate` skill

## Rule 4 — Persistence axis parity

Each `internal/adapter/persistence/<dialect>/` folder (postgres,
mysql, mongodb) MUST compile + satisfy the same domain `Repository`
port. When you rename a domain entity, change a port signature, or
add a new port method:

- Update all three dialect adapters.
- Update migrations: `db/migrations/<dialect>/NNN_<name>.up.sql`.
- Update seeds: `db/seeds/<dialect>/NNN_<name>.{sql,json}`.

Only postgres has full testcontainers integration tests today; mysql
and mongodb are compile-parity only. Both still build clean via
`go build ./...`.

→ Adapter catalog: `docs/ADAPTERS.md`
→ Migration policy: `db/migrations/README.md`
→ Seed policy: `db/seeds/README.md`

## Rule 5 — Migration & seed policy

- **Forward-only.** No `.down.sql` files in this repo. Local-dev
  rollback uses drop + remigrate.
- File prefix `NNN_` (zero-padded, monotonic per dialect):
  `001_create_users.up.sql`, `002_create_todos.up.sql`.
- Seeds are idempotent (`ON CONFLICT DO NOTHING` for postgres,
  `INSERT IGNORE` for mysql, `$setOnInsert` upsert for mongodb).

## Rule 6 — Auth axis interface contract

`internal/adapter/auth/<option>/` (jwt, apikey, oauth2-oidc, none)
implement the same middleware port. A new auth option MUST:

- Satisfy the domain `auth` port interface.
- Provide an echo middleware factory.
- Compile + satisfy the interface even when stubbed. Stub options
  (`oauth2-oidc`, `none`) may throw "not implemented" at runtime but
  the type contract must hold.

→ Detail: `docs/ADAPTERS.md#auth`

## Rule 7 — Testing pattern

| Layer | Path | Tooling | Coverage target |
|---|---|---|---|
| Unit (domain) | `internal/domain/<x>/*_test.go` | testify + handwritten mocks | ≥ 80% |
| Handler | `internal/http/handler/*_test.go` | httptest + mock service | ≥ 70% |
| Integration (adapter) | `internal/adapter/persistence/<x>/*_test.go` (build tag `integration`) | testcontainers-go | ≥ 60% per adapter |

Test name encodes FR: `TestList_FR001_ReturnsTodosForUser`. Run via
`make test-unit`, `make test-integration`, `make test-cover`.

→ Reference tests: `internal/domain/todo/service_test.go`,
  `internal/adapter/persistence/postgres/todo_repo_test.go`

## Rule 8 — Config & secrets

- Typed config in `internal/config/config.go`, loaded via Viper from
  `env/env.<APP_ENV>` (defaults to `local`).
- Per-axis-option config fields wrapped in
  `boilerplate:axis=X option=Y START..END` markers.
- Secret resolution via `internal/adapter/secrets/<option>/` (env /
  gcp / aws / vault). The `env` option is implicit via `config.Load`;
  cloud sources are wired manually in `container.go`.
- NEVER log secret values. Use structured logger redaction.
- **env/env.example parity:** Every env var read in `config.go` MUST
  have a matching entry in `env/env.example` with a descriptive comment
  and an empty or safe default value (never real credentials). Add the
  entry in the same commit as the config field. Missing entries are a
  review gate enforced by the lintas `sdlc-be-go:review-implementation`
  skill.

## Rule 9 — Logging

- zerolog structured. Get the logger via `pkg/logger`.
- Levels:
  - `debug` — verbose dev tracing (not in prod by default)
  - `info` — request flow, lifecycle events
  - `warn` — recoverable issue (deprecated path, retry)
  - `error` — 5xx-class failure
- Always include `requestID` from `pkg/ctxutil`.

## Rule 10 — Commit & PR

- Conventional Commits: `<type>(<scope>): <subject>` with body
  explaining **why**, not **what**.
- `feat` and `fix` REQUIRE `Refs: <TICKET-ID>` footer.
- Branch: `<type>/<TICKET-ID>-<slug>` (≤ 50 chars, kebab-case slug).
- NEVER `--no-verify`. Pre-push runs `go test ./...` + `golangci-lint run`.

→ Hook detail: `docs/CONTRIBUTING.md`

## Reference domain

`internal/domain/todo/` is the end-to-end reference implementation.
Use it as the pattern when adding a real domain.

Entity (see `internal/domain/todo/todo.go`):

```go
type Todo struct {
    ID          uuid.UUID
    Title       string
    Description string
    DueDate     *time.Time
    CreatedAt   time.Time
    CreatedBy   uuid.UUID  // ownership marker; List filters on this
    ModifiedAt  time.Time
    ModifiedBy  uuid.UUID
}
```

Endpoints exposed by default:

| Method | Path | Description |
|---|---|---|
| GET | `/healthz` | Liveness probe (no auth, no envelope) |
| POST | `/v1/auth/login` | Issue access + refresh tokens |
| POST | `/v1/auth/logout` | Revoke session |
| POST | `/v1/auth/renew` | Rotate access token |
| POST | `/v1/todos` | Create todo — `201` + `SuccessEnvelope` |
| GET | `/v1/todos` | List todos — paginated, `PaginatedEnvelope` |
| GET | `/v1/todos/:id` | Get todo by ID |
| PUT | `/v1/todos/:id` | Update todo |
| DELETE | `/v1/todos/:id` | Soft delete — `200` + `data: {}` |

## rpc axis (gRPC)

The `rpc` axis (`grpc` / `none`, default `none`) adds gRPC transport support.

**Prerequisite:** install buf — `brew install bufbuild/buf/buf`

**Key paths:**
- `proto/` — `.proto` source files
- `gen/` — committed generated stubs (`*.pb.go`, `*_grpc.pb.go`)
- `internal/adapter/rpc/grpc/` — server bootstrap, client factory, health service
- `buf.yaml`, `buf.gen.yaml` — buf module + codegen config

**Regenerate stubs:** `make buf-generate` (or `buf generate`)

Generated stubs are committed so `go build` works without requiring buf at build
time. Re-run `buf generate` only when `.proto` files change.

**Marker discipline:** all gRPC wiring (config, imports, container construction,
server start/stop) is wrapped in `// boilerplate:axis=rpc option=grpc START/END`
markers so `rpc=none` prune removes them cleanly.

## inference axis

The `inference` axis (`noop` / `onnxruntime`, default `noop`) wires a local ML
`Engine` (`internal/domain/inference/`) into the container via
`internal/app/inference.go` (`newInference`, switches on `INFERENCE_BACKEND`).
`noop` returns empty results; `onnxruntime` is a compile-clean stub until a
product wires real models (`docs/ML_INFERENCE_PATTERN.md`). All wiring is
marker-wrapped `// boilerplate:axis=inference option=<x> START/END`.

## Quick reference

| Topic | Source |
|---|---|
| Layering rules + diagrams | `docs/ARCHITECTURE.md` |
| Axis catalog (45 options) | `docs/ADAPTERS.md` |
| Error code catalog | `docs/ERROR_CODES.md` |
| Pruning algorithm | `docs/PRUNING.md` |
| Commit/PR/hook detail | `docs/CONTRIBUTING.md` |
| Migration policy | `db/migrations/README.md` |
| Seed policy | `db/seeds/README.md` |
| Manifest | `.boilerplate.yaml` |

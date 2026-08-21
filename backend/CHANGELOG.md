# Changelog

All notable changes to `backend-belajar-lintas-ronald` boilerplate are documented here.
The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added — TodoApp minimum runnable

- **Auth domain** (`internal/domain/auth/`): `User` entity, `Authenticator` service, `UserRepository` / `SessionStore` / `TokenIssuer` ports.
- **JWT issuer** (`internal/adapter/token/jwt/`): HS256 signing/parsing with configurable access + refresh TTL.
- **Redis session store** (`internal/adapter/session/redis/`): refresh-token storage with `refresh:<hash>` keys and TTL.
- **User repositories**: postgres (integration-tested via testcontainers), mysql, mongodb (compile parity).
- **Auth handler** (`internal/http/handler/auth_handler.go`): `POST /v1/auth/{login,logout,renew}`.
- **Migrations + seeds** (3 dialects): `001_create_users` + `002_create_todos`. Demo accounts (`demo1@example.com`, `demo2@example.com` — password `Demo1234!`, bcrypt cost 12) and 5 pre-seeded todos.
- **Bcrypt hash helper** (`scripts/gen-seed-hash.go`).
- **JWT env config**: `JWT_SECRET`, `JWT_ACCESS_TTL_SEC`, `JWT_REFRESH_TTL_SEC` (root `compose.yml` `backend-belajar-lintas-ronald` service, marker-wrapped `axis=auth option=jwt`).
- **`data/.gitkeep`**: satisfies pre-existing `COPY ./data ./data` reference in `DockerfileLocal`.
- **`customerror.CustomError.Is`**: `errors.Is(enriched, sentinel)` now returns `true` when both share `code`/`httpCode`/`message` (metadata-insensitive — same semantics as the existing `Equal` method).

### Changed (BREAKING — domain rename)

- **Domain `item` renamed to `todo`** with adjusted entity:
  - `ID`: `string` → `uuid.UUID`
  - Drop `OwnerID` + `Status` fields
  - Add `Description`, `DueDate`, `CreatedBy`, `ModifiedBy` (audit columns)
  - Rename `UpdatedAt` → `ModifiedAt`
- **Handler routes** `/v1/items/*` → `/v1/todos/*` (paginated via `?page=&page_size=`).
- **Migrations** renumbered: users is `001_`, todos is `002_`.
- **OpenAPI spec** rewritten to include `/v1/auth/*` + `/v1/todos/*` paths with envelope schemas.
- **DockerfileLocal** Go base image: `golang:1.24.3-bookworm` → `golang:1.25-bookworm` (matches `go.mod` `go 1.25.0`).

Plugin SDLC `sdlc-be-go` in `lintas` repo must be updated to handle the new auth axis routes + `data/` folder placement (separate follow-up).

### Changed (BREAKING — orchestration)

- **Docker Compose moved to monorepo root.** `golang/compose.yml` has been deleted; the boilerplate now uses a single `compose.yml` at the repo root that covers every subproject's services.
- **`golang/scripts/prune.sh` and `golang/scripts/verify-prune.sh` moved to root `scripts/`.** New flag: `--subproject=<name>` (required) plus `--selections=<axis=opt,…>`. The pruner now parses a `subproject=<list>` marker dimension.
- **Marker syntax extended**: blocks now carry `subproject=<comma-list>` in addition to the existing `axis=…` / `option=…` fields. Previous syntax (`axis=… option=…` only) is no longer recognised.
- **`.boilerplate.yaml` axis `container`** no longer declares a `compose` option; compose is monorepo-level concern.
- **Local dev** commands change: invoke `docker compose -f compose.yml <cmd>` from the repo root, not from `golang/`.

### Added (orchestration)

- `frontend-belajar-lintas-ronald` service definition in root `compose.yml` so the FE side can be containerized via `docker compose up -d frontend-belajar-lintas-ronald`.
- Root `compose.yml` markers now support nested blocks (outer `subproject=` wrapping inner `axis=`/`option=` blocks) via a stack-based parser in `scripts/prune.sh`.

### Added

- `db:seed` Cobra subcommand with per-dialect dispatch for Postgres, MySQL, and MongoDB.
- Sample seed files at `db/seeds/{postgres,mysql,mongodb}/001_items.{sql,json}` demonstrating the idempotency patterns.
- `db/seeds/README.md` documenting seed conventions and `db/migrations/README.md` documenting the forward-only migration policy.
- Makefile targets `seed-postgres`, `seed-mysql`, `seed-mongodb` (markered for plugin pruning).
- Internal config fields `Config.MySQL` and `Config.Mongo` (markered) so `db:seed` can connect to the selected dialect.

### Removed

- `db/migrations/{postgres,mysql}/001_create_items.down.sql` — boilerplate is now forward-only. Local-dev rollback uses drop + remigrate.

### Changed (BREAKING — HTTP contract)

- **Response envelope** now follows the lintas API contract. Every endpoint returns `{success, code, message, data, timestamp}` for success (with `pagination` for lists) and `{success, code, message, timestamp, metadata?}` for errors. The previous shape (`{error, error_description, timestamp, metadata}`) is gone.
- **Error `code`** is a 5-digit numeric (`HTTP_PREFIX + 2-digit subcode`, e.g., `40400` for not-found). Replaces the previous string codes (`"item_not_found"`, etc.).
- **`pkg/customerror.ErrorCode`** type changed from `string` to `int`. All callers that asserted string codes need to use the new numeric constants in `pkg/customerror/error_code.go`.
- **`DELETE /v1/items/:id`** now returns `200 OK` + envelope `data: {}` instead of `204 No Content` (envelope universality).
- **List query params** renamed `limit/offset` → `page/page_size` to match the response `pagination` block end-to-end.
- **Validation failures** now return HTTP `422 Unprocessable Entity` (code `42200`) instead of `400 Bad Request`. Body parse failures still return `400` (code `40000`).

### Added

- `pkg/customerror` generic catalog with 8 pre-built codes (`CodeBadRequest=40000`, `CodeUnauthorized=40100`, ..., `CodeInternalServer=50000`) and matching `Err*` vars.
- `pkg/customerror.NewConflictError` constructor (HTTP 409).
- `pkg/pagination.Page` struct matching the lintas pagination shape, with `Build(page, pageSize, totalData)` helper.
- `internal/http/response` helpers: `OK`, `Paged`, `Err` (replace per-status helpers).
- `docs/ERROR_CODES.md` documenting the 5-digit numeric code convention and service-extension guidelines.
- `docs/ARCHITECTURE.md` section on response envelope and package boundaries.

### Removed

- `pkg/customerror/error_code.go` legacy auth-specific constants (~200 entries, e.g., `ErrPINMustBeNumeric`, `InvalidOTP`) — not boilerplate concerns.
- `internal/http/response` per-status helpers (`BadRequest`, `Unauthorized`, `Forbidden`, `NotFound`, `Conflict`, `InternalError`) — collapsed into the single `Err(c, err)` helper.
- `internal/http/handler/dto.ListItemsResponse` — envelope owns the pagination shape.

## [v1.0.0] — 2026-05-13

### Added

- **Hexagonal Architecture (Ports & Adapters)** layout:
  `internal/{domain,adapter,http,app,config}` + `pkg/`.
- **Reference `item` CRUD domain** demonstrating all hexagonal patterns:
  entity, ports (Repository, Cache), service with cache-aside +
  soft-delete, mocks, unit tests, handler tests, integration tests.
- **10 modular axes with 34 adapter options:**
  - persistence: postgres, mysql, mongodb
  - cache: redis, inmemory, noop
  - auth: jwt, apikey, oauth2-oidc, none
  - observability: otel, datadog, noop
  - messaging: gcp-pubsub, rabbitmq, redis-streams, aws-sqs, aws-sns, noop
  - storage: gcs, s3, minio, local, noop
  - secrets: env, gcp, aws, vault
  - notification: smtp, sendgrid, noop
  - api-doc: swag, oapi-codegen, none
  - jobs: cron, worker, noop
  - container: dockerfile, compose, tilt, nginx
- **Machine-readable manifest** at `.boilerplate.yaml` (schema version 1).
- **Comment markers** (`boilerplate:axis=X option=Y START..END`) in
  wiring files for programmatic pruning by the SDLC plugin.
- **`scripts/prune.sh`** — interactive manual pruner.
- **`scripts/verify-prune.sh`** — post-prune sanity check.
- **3-layer commit convention enforcement:**
  - commit-msg hook (Conventional Commits + Refs footer)
  - pre-push hook (branch name validation)
  - GitHub Action (PR title validation)
- **Documentation:** README, ARCHITECTURE, ADAPTERS, PRUNING, CONTRIBUTING.
- **Tooling:** mise + `.tool-versions`, Makefile with `go run @version`,
  golangci-lint v2 with hexagonal depguard rules.
- **Per-dialect SQL migrations** at `db/migrations/<dialect>/`.

### Changed

- Module path: `sqeid-api` → `backend-belajar-lintas-ronald`.
- Env var: `SQEID_ENV` → `APP_ENV`.

### Removed

- Sqeid-specific business code (Keycloak/OTP/PII/Realm handlers, services,
  repositories).
- Nix shells (`shell.nix`, `shell-wsl.nix`) — replaced by `.tool-versions`.

### Notes

- Boilerplate ships with all axes coexisting (~34 adapter folders) so the
  SDLC plugin can prune to any combination. Default selection used for
  the `item` reference: postgres + redis + jwt + otel.

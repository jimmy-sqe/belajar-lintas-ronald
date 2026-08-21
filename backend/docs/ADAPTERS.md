# Adapters Catalog

The boilerplate ships **14 modular axes** with **45 adapter options**. Each option is a separate folder under `internal/adapter/<axis>/<option>/`. The SDLC plugin (or `scripts/prune.sh`) prunes unselected options.

The canonical schema is `.boilerplate.yaml` at the repo root. This document explains when to choose each option.

## Multiplicity rules

- **`exactly_one`** — must select exactly one option. A `noop`/`none` option is provided as the explicit "off" choice.
- **`one_or_more`** — select at least `min` options; combinations allowed.

## Integration modes

Each axis declares an `integration_mode` in `.boilerplate.yaml` so automated
consumers (plugins, scaffolders) can tell which selections need manual wiring
after pruning.

- **`wired`** — adapter is instantiated in `internal/app/container.go`; the
  default selection works out-of-the-box after pruning.
- **`standalone`** — adapter folder ships with example code, but is **not**
  referenced from `container.go`/`config.go`. The user wires it manually when
  adopting the capability.
- **`meta`** — affects build/tooling only; no Go runtime wiring.

| Axis | Mode | Notes |
|---|---|---|
| `persistence` | wired | DB pool created in container; repository injected into services. |
| `cache` | wired | Cache client created in container; injected into services that need it. |
| `auth` | wired | Auth middleware mounted on `/v1/*` routes. |
| `observability` | wired | Tracer/meter providers initialised at startup. |
| `secrets` | wired | `env` option is implicit via `config.Load`. `gcp`/`aws`/`vault` are standalone — wire manually in `container.go` when adopting an external secret manager. |
| `messaging` | standalone | Adapters ship as references. Wire publishers/consumers in `container.go` when needed. |
| `storage` | wired | Object-store client selected at startup via `STORAGE_BACKEND` env var. Supports multi-env setups (e.g., minio locally, GCS in production). |
| `notification` | standalone | SMTP/Sendgrid clients ship as references. Wire into the service that sends notifications. |
| `jobs` | standalone | Cron/worker entry points ship as references. Mount in `cmd/` or `container.go` per your scheduling needs. |
| `api-doc` | meta | Affects Makefile targets and generated docs only. |
| `container` | meta | Affects Dockerfile/compose/Tilt/nginx files; no Go wiring. |

## persistence — `one_or_more` (min 1)

Required by the `todo` reference domain (Repository port).

| Option | When to choose |
|---|---|
| `postgres` | Default. Strong consistency, mature ecosystem, JSON support, full SQL feature set. |
| `mysql` | Compatibility with existing MySQL infra (e.g., AWS Aurora MySQL, PlanetScale). |
| `mongodb` | Flexible schema, document-oriented data, log-like workloads. Choose alongside `postgres`/`mysql` if you have both relational and document needs. |

Migrations live at `db/migrations/<dialect>/` (golang-migrate format).

### Seeding

Each persistence dialect ships with sample seed files under `db/seeds/<dialect>/`:

- **Postgres / MySQL** — `NNN_<name>.sql` files executed in lexicographic order via `db.ExecContext`. Idempotency is per-file (`ON CONFLICT DO NOTHING` / `INSERT IGNORE`).
- **MongoDB** — `NNN_<collection>.json` arrays of documents. The executor upserts each document via `$setOnInsert` so re-runs do not overwrite existing rows.

Run seeds: `make seed-<dialect>` (which wraps `go run . db:seed`). The command auto-routes to the dialect declared in config.

See `db/seeds/README.md` for the full rules. See `db/migrations/README.md` for the forward-only migration policy.

## cache — `exactly_one`

Optional capability used by `todo.Service` for read-through caching.

| Option | When to choose |
|---|---|
| `redis` | Default. Shared cache across pods, persistent if needed. |
| `inmemory` | Single-instance services or tests; no cross-pod coherence. |
| `noop` | Don't want caching; satisfies the interface so service code doesn't change. |

## auth — `exactly_one`

Authentication middleware applied to `/v1/*` routes.

| Option | When to choose |
|---|---|
| `jwt` | Default. HS256 self-issued tokens or symmetric keys. |
| `apikey` | Service-to-service or simple API key authentication. |
| `oauth2-oidc` | OIDC-compliant IdP (Keycloak, Auth0, Cognito, Okta, Google, …). Validates via JWKS. |
| `none` | Auth happens at the edge (API gateway), or service is internal-only. |

## observability — `exactly_one`

Tracing + structured logging glue.

| Option | When to choose |
|---|---|
| `otel` | Default. Vendor-neutral; export to any OTLP-compatible backend (Tempo, Honeycomb, …). |
| `datadog` | Datadog APM. |
| `noop` | Local development without an APM target. |

## messaging — `one_or_more` (min 1)

Standalone publishers for outbound events. Not integrated with `todo` by default — wire into your service constructor when you need to emit events.

| Option | When to choose |
|---|---|
| `gcp-pubsub` | GCP-hosted services. |
| `rabbitmq` | Self-hosted broker; supports topic exchanges, routing keys. |
| `redis-streams` | Lightweight broker; reuses `cache.redis` infra. |
| `aws-sqs` | Queue (point-to-point). |
| `aws-sns` | Pub/sub (fan-out). Common pair with SQS. |
| `noop` | No outbound events. |

## storage — `one_or_more` (min 1)

Object storage for file uploads, exports, attachments. Wired at startup via
`STORAGE_BACKEND` env var — supports per-environment adapter switching.

| Option | When to choose |
|---|---|
| `gcs` | GCP. |
| `s3` | AWS or any S3-compatible service. |
| `minio` | Self-hosted S3-compatible. |
| `local` | Local filesystem (dev or single-instance services). |
| `noop` | No file storage. |

### Multi-environment setup

Select multiple options at scaffold time (e.g., `[minio, gcs]`) so both adapter
folders survive pruning. Then set `STORAGE_BACKEND` per environment:

| Env file | `STORAGE_BACKEND` | Why |
|---|---|---|
| `env.local` | `minio` | Free, self-hosted, no cloud credentials needed |
| `env.staging` | `gcs` | Match production infra |
| `env.production` | `gcs` | Google Cloud Storage |

The service reads `STORAGE_BACKEND` at startup and instantiates the matching
adapter. If the value is empty or unrecognised, startup fails with an explicit
error listing valid options.

Single-option selection still works — pick one option and set `STORAGE_BACKEND`
to that value in all env files.

## secrets — `one_or_more` (min 1)

Secret-source adapters. `env` is always available; cloud sources are additive.

| Option | When to choose |
|---|---|
| `env` | Default; reads from process env. |
| `gcp` | Google Secret Manager. |
| `aws` | AWS Secrets Manager. |
| `vault` | HashiCorp Vault. |

## notification — `exactly_one`

Outbound email.

| Option | When to choose |
|---|---|
| `smtp` | Generic SMTP (any provider). |
| `sendgrid` | Twilio SendGrid. |
| `noop` | No email. |

## api-doc — `exactly_one`

API documentation strategy.

| Option | When to choose |
|---|---|
| `swag` | Code-first via inline comments on handlers. Generates `swagger.json`. |
| `oapi-codegen` | Spec-first; types and Echo server stubs generated from `openapi.yaml`. |
| `none` | No API docs (or maintain externally). |

## jobs — `exactly_one`

In-process background work.

| Option | When to choose |
|---|---|
| `cron` | Scheduled tasks via `robfig/cron/v3`. |
| `worker` | Channel-driven worker pool. Pair with a `messaging` adapter for queue consumers. |
| `noop` | No background jobs. |

## rpc — `exactly_one`

gRPC transport layer (alongside the always-present Echo HTTP server).

| Option | When to choose |
|---|---|
| `none` | Default. HTTP-only service. |
| `grpc` | Adds a gRPC server bound to `GRPC_HOST:GRPC_PORT`. Stubs are committed under `gen/` (regenerate with `make buf-generate`), so `go build` needs no `buf` at build time. |

**Key paths:**
- `proto/` — `.proto` source files
- `gen/` — committed generated stubs (`*.pb.go`, `*_grpc.pb.go`)
- `internal/adapter/rpc/grpc/` — server bootstrap, client factory, health service

All wiring (config, container construction, server start/stop) is marker-wrapped
`// boilerplate:axis=rpc option=grpc` so an `rpc=none` prune removes it cleanly.

## inference — `exactly_one`

Local ML model inference port. Follows the same hexagonal switch-factory pattern
as persistence and cache. See `docs/ML_INFERENCE_PATTERN.md` for the full pattern.

| Option | When to choose |
|---|---|
| `noop` | Default. Returns an empty `Result`. No ML dependencies. |
| `onnxruntime` | ONNX Runtime engine reading models from `INFERENCE_MODEL_DIR`. Shipped as a compile-clean stub; the product adds the real `github.com/yalue/onnxruntime_go` binding via `go get` per `docs/ML_INFERENCE_PATTERN.md`. |

Wired in `internal/app/inference.go` (`newInference`), which switches on
`INFERENCE_BACKEND`.

## container — `one_or_more` (min 1)

Local dev / packaging files.

| Option | When to choose |
|---|---|
| `dockerfile` | Always. Production image build. |
| `compose` | Local dev orchestration. |
| `tilt` | Multi-service hot-reload during dev. |
| `nginx` | Reverse-proxy front for local (e.g., to mimic prod hostnames). |

## sample-app — `one_or_more`

Reference domain shipped as an opinionated end-to-end example. Removable.

| Option | When to choose |
|---|---|
| `todo-app` | Default. Keeps `internal/domain/todo/`, the todo HTTP handler + DTOs, todo persistence/cache adapters across all dialects, and related tests. Remove when building a real product domain. |

## Adding a new option to an existing axis

See `docs/CONTRIBUTING.md`.

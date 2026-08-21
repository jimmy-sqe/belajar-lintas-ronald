# Go Service Boilerplate (`backend-belajar-lintas-ronald`)

Greenfield Go service template using **Hexagonal Architecture (Ports & Adapters)** with a **machine-readable modular adapter catalog**. Ships with the lintas team's response envelope, numeric error-code scheme, and an end-to-end reference domain so a new service is runnable on day 1.

> Part of the [boilerplate-monorepo](../README.md). See the monorepo root for the cross-language overview.

The template ships with **all 11 modular axes** (34 adapter options) coexisting. The SDLC plugin (`sdlc-be-go`) — or a human via [`scripts/prune.sh`](scripts/prune.sh) — selects the combination a real service keeps; the rest is removed automatically.

---

## Contents

- [Who is this for](#who-is-this-for)
- [Quick start](#quick-start)
- [Architecture at a glance](#architecture-at-a-glance)
- [Modular axes](#modular-axes)
- [Reference domain (`todo`)](#reference-domain-todo)
- [Response envelope & error codes](#response-envelope--error-codes)
- [Tooling](#tooling)
- [Commit & PR conventions](#commit--pr-conventions)
- [Documentation](#documentation)
- [License](#license)

---

## Who is this for

| You are… | Use this template via… |
|---|---|
| Bootstrapping a **new Go service** in S-Quantum Engine | The SDLC plugin: run `sdlc-be-go:scaffold-be-project <project-name>` — it clones a tagged release of this folder, prunes unselected axes, and resets git history. |
| Adding **a new axis or option** to the template itself | This repo directly. Read [`docs/PRUNING.md`](docs/PRUNING.md) for marker syntax and [`docs/CONTRIBUTING.md`](docs/CONTRIBUTING.md) for the contribution workflow. |
| Just experimenting | Clone, run [`scripts/prune.sh`](scripts/prune.sh) interactively to keep only the axes you need. |

---

## Quick start

### Path A — Via the SDLC plugin (recommended)

```bash
# From a workspace where the lintas plugin pack is installed
sdlc-be-go scaffold-be-project my-service
cd my-service
make run
```

The plugin reads `.boilerplate.yaml`, asks which option to keep per axis, removes the rest, and produces a fresh repo with a clean commit history.

### Path B — Manual, for local exploration

```bash
# 1. Clone this monorepo
git clone <this-repo-url>
cd boilerplate-monorepo

# 2. Install per-project tools (Go 1.25, Tilt, …) via mise
(cd golang && mise install)

# 3. (Optional) Preview a prune for a target subproject
bash scripts/prune.sh --subproject=golang \
    --selections=persistence=postgres,cache=redis,storage=gcs \
    --output=/tmp/preview.yml

# 4. Bring up infra + service via the root compose
docker compose up -d postgres redis backend-belajar-lintas-ronald

# 5. Verify the service is up
curl http://localhost:8000/healthz
# → {"status":"ok"}
```

All `docker compose` commands run from the repo root, not from `golang/`. The root `compose.yml` carries `subproject=` markers so the SDLC plugin can extract a Go-only subset when scaffolding new services.

---

## Architecture at a glance

```
internal/
├── domain/<feature>/         ← pure business logic; ports declared as interfaces
├── adapter/<axis>/<option>/  ← infrastructure adapters that implement ports
├── http/                     ← inbound HTTP layer (Echo) + envelope renderer
├── app/                      ← composition root: container.go wires everything
└── config/                   ← Viper-loaded typed config

pkg/                          ← shared utilities (logger, ctxutil, customerror, pagination, …)
cmd/                          ← Cobra subcommands (serve, db:migrate, db:seed)
db/migrations/<dialect>/      ← forward-only SQL migrations per persistence backend
db/seeds/<dialect>/           ← idempotent initial/master data per persistence backend
.boilerplate.yaml             ← machine-readable manifest of axes + options
scripts/prune.sh              ← interactive pruner for axis options
```

**Start reading at**: [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) (layering rules, dependency direction, response envelope) → [`docs/ADAPTERS.md`](docs/ADAPTERS.md) (full axis catalog, when-to-choose guidance) → [`docs/PRUNING.md`](docs/PRUNING.md) (manifest schema, marker syntax, prune algorithm).

---

## Modular axes

14 axes / 45 options. Default selection used by the reference `todo` domain is in **bold**.

| Axis | Multiplicity | Options |
|---|---|---|
| persistence | one_or_more (min 1) | **postgres**, mysql, mongodb |
| cache | exactly_one | **redis**, inmemory, noop |
| auth | exactly_one | **jwt**, apikey, oauth2-oidc, none |
| observability | exactly_one | **otel**, datadog, noop |
| messaging | one_or_more (min 1) | **gcp-pubsub**, rabbitmq, redis-streams, aws-sqs, aws-sns, noop |
| storage | exactly_one | **gcs**, s3, minio, local, noop |
| secrets | one_or_more (min 1) | **env**, gcp, aws, vault |
| notification | exactly_one | **smtp**, sendgrid, noop |
| api-doc | exactly_one | **swag**, oapi-codegen, none |
| jobs | exactly_one | **cron**, worker, noop |
| container | one_or_more (min 1) | **dockerfile**, tilt, nginx |
| rpc | exactly_one | **none**, grpc |
| inference | exactly_one | **noop**, onnxruntime |
| sample-app | one_or_more | **todo-app** |

Full schema: [`.boilerplate.yaml`](.boilerplate.yaml). Pruning algorithm: [`docs/PRUNING.md`](docs/PRUNING.md).

---

## inference axis

The `inference` axis (`noop` / `onnxruntime`, default `noop`) wires a local ML
inference `Engine` into the container.

- Port: `internal/domain/inference/engine.go` (`Engine.Infer`).
- Adapters: `internal/adapter/inference/noop` (empty results, the default) and
  `internal/adapter/inference/onnxruntime` (a compile-clean stub — products
  uncomment the real `github.com/yalue/onnxruntime_go` binding and `go get` it
  per `docs/ML_INFERENCE_PATTERN.md`).
- Wiring: `internal/app/inference.go` `newInference` switches on
  `INFERENCE_BACKEND`; `container.go` stores the `Engine` on the container for
  products to consume.

Marker discipline: all inference wiring is wrapped in
`// boilerplate:axis=inference option=<noop|onnxruntime> START/END`.

---

## Reference domain (`todo`)

A complete CRUD reference at [`internal/domain/todo/`](internal/domain/todo/) plus its HTTP handlers, Postgres adapter (with golang-migrate migrations), Redis cache, JWT auth middleware, and OTel observability — all wired as the default selection. Use it as the pattern when adding real domains.

Endpoints exposed by default:

| Method | Path | Description |
|---|---|---|
| GET | `/healthz` | Liveness probe (no auth, no envelope) |
| POST | `/v1/todos` | Create todo — returns `201` + `SuccessEnvelope` |
| GET | `/v1/todos` | List todos — paginated via `?page=&page_size=`, returns `PaginatedEnvelope` |
| GET | `/v1/todos/:id` | Get todo by ID |
| PUT | `/v1/todos/:id` | Update todo |
| DELETE | `/v1/todos/:id` | Soft delete — returns `200` + `SuccessEnvelope` with `data: {}` |

---

## Response envelope & error codes

All responses (except `/healthz`) follow the **lintas team API contract**:

- **Success**: `{ success: true, code: 200, message, data, timestamp }` (lists add `pagination`)
- **Error**: `{ success: false, code: 40400, message, timestamp, metadata? }` — `code` is a 5-digit numeric (`HTTP_PREFIX + 2-digit subcode`)

The boilerplate ships 8 generic error codes (`40000`, `40100`, `40300`, `40400`, `40900`, `42200`, `42900`, `50000`); scaffolded services extend the catalog with domain-specific subcodes in the `01–99` range.

→ Full convention: [`docs/ERROR_CODES.md`](docs/ERROR_CODES.md). Envelope contract: [`docs/ARCHITECTURE.md#response-envelope`](docs/ARCHITECTURE.md#response-envelope).

---

## Tooling

**Runtime**

- **Go 1.25** — pinned in `go.mod` + `.tool-versions` (mise/asdf)
- **Echo v4** — HTTP router & middleware
- **sqlx + Squirrel** — SQL access & query builder
- **Viper** — typed config loaded from `env/env.<APP_ENV>`
- **zerolog** — structured logging

**Dev & test**

- **testify + testcontainers-go** — unit + integration tests against real containers
- **golangci-lint, mockery, golang-migrate** — invoked via `go run @version` in the Makefile (no global install required)
- **mise** — toolchain pin (`.tool-versions`)

**Make targets**

```bash
make help               # list every target
make test               # unit tests
make test-integration   # integration tests (build tag, runs testcontainers)
make test-cover         # coverage report
make lint               # golangci-lint
make build              # binary at bin/$(APP_NAME)
make run                # start service locally
make migrate-postgres   # apply migrations (per dialect: -mysql, -mongodb)
make seed-postgres      # apply seeds (per dialect: -mysql, -mongodb)
```

---

## Commit & PR conventions

This repo uses **Conventional Commits** with a mandatory `Refs: <TICKET-ID>` footer for `feat`/`fix` types. Branch names follow `<type>/<TICKET-ID>-<slug>` (types: `feature | fix | chore | docs | experiment`).

Enforced by a 3-layer setup:

1. **commit-msg hook** — Conventional Commits format + `Refs:` footer check
2. **pre-push hook** — branch-name validation + `go test ./...` + `golangci-lint run`
3. **GitHub Action** — PR title validation on open/edit

Never use `--no-verify` to bypass hooks. See [`docs/CONTRIBUTING.md`](docs/CONTRIBUTING.md).

---

## Documentation

| Doc | What's in it |
|---|---|
| [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) | Hexagonal layering rules, dependency direction, cross-domain access pattern, response-envelope contract |
| [`docs/ADAPTERS.md`](docs/ADAPTERS.md) | Full axis catalog with when-to-choose guidance per option |
| [`docs/ERROR_CODES.md`](docs/ERROR_CODES.md) | 5-digit numeric error-code convention + service-extension guidelines |
| [`docs/PRUNING.md`](docs/PRUNING.md) | Manifest schema, marker syntax, prune algorithm |
| [`docs/CONTRIBUTING.md`](docs/CONTRIBUTING.md) | Branch/commit/PR conventions, how to add new axis options |
| [`db/migrations/README.md`](db/migrations/README.md) | Forward-only migration policy |
| [`db/seeds/README.md`](db/seeds/README.md) | Idempotent seed conventions per dialect |

For original design history (specs + implementation plans), see [`../docs/superpowers/`](../docs/superpowers/) at the monorepo root.

---

## License

See [`LICENSE`](../LICENSE) at the monorepo root (MIT by default).

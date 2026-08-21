# Boilerplate Monorepo — S-Quantum Engine Service Templates

Greenfield service boilerplates for the S-Quantum Engine SDLC, organised by language under sibling subfolders. Each subfolder is a **self-contained boilerplate** with its own manifest (`.boilerplate.yaml`), documentation, and tooling.

New services are scaffolded from these templates by the SDLC plugin pack (`sdlc-be-go`, `sdlc-fe-nextjs`, …). The plugin clones the relevant subfolder at a tagged version and prunes axes the team didn't pick.

## Subprojects

| Path | Stack | Status | Docs |
|---|---|---|---|
| [`backend/`](backend/) | Go 1.25 + Echo v4 + sqlx + Squirrel + Redis + hexagonal (ports & adapters) | v1.0.0 | [README](backend/README.md) · [manifest](backend/.boilerplate.yaml) |
| [`frontend/`](frontend/) | Next.js 15 + React 19 + TS 5 + Tailwind + Horizon UI + Kubb | v1.0.0 | [README](frontend/README.md) · [manifest](frontend/.boilerplate.yaml) |
| [`qa-test/`](qa-test/) | Python 3.12 + pytest + Playwright (UI) + requests/faker (API), single-project | v1.0.0 | [README](qa-test/README.md) · [manifest](qa-test/.boilerplate.yaml) |
| `python/` | (FastAPI — TODO) | — | — |
| [`java/`](java/) | Java 25 + Spring Boot 4.1 + JPA + Flyway + Caffeine + hexagonal (ports & adapters) | v0.1.0 | [README](java/README.md) · [manifest](java/.boilerplate.yaml) |

## Cross-subproject axes

Some axes are shared across multiple subprojects (Go, Python). See each subproject's `docs/ADAPTERS.md` for the full axis catalog.

| Axis | Options | Description |
|---|---|---|
| `rpc` | `grpc` · `none` | gRPC transport (server + client). Requires buf. Default `none`. |

## Layout

```
.
├── compose.yml            Root Docker Compose: all subprojects' services with subproject= markers
├── scripts/               Monorepo-level tooling (prune.sh, verify-prune.sh)
├── backend/                Go service boilerplate (hexagonal, 11 modular axes)
├── frontend/                Next.js boilerplate (App Router, 10 modular axes)
├── qa-test/               QA automation boilerplate (pytest + Playwright, 4 axes)
├── docs/superpowers/      Monorepo design specs + implementation plans (history)
├── .github/workflows/     Shared monorepo CI (path-scoped per subproject)
├── .pre-commit-config.yaml Pre-commit hooks (Go scoped via files: ^backend/)
├── .gitignore             Monorepo-wide ignore rules
├── .gitleaks.toml         Secret scanning config (whole repo)
└── .gitmessage            Commit message template (Conventional Commits)
```

Per-language `.gitignore`, `.tool-versions`, `Dockerfile`, etc. live inside each subproject.

## Quick start — Tilt (hot-reload dev, recommended)

[Tilt](https://docs.tilt.dev/install.html) gives live source sync without rebuilding Docker images on every change:
- **Go backend**: file save → sync + binary rebuild → ~2–5 s
- **Next.js frontend**: file save → Turbopack HMR → browser update <1 s

### Install Tilt

```bash
# macOS and Linux
curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash

# Windows (PowerShell)
iex ((new-object net.webclient).DownloadString('https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.ps1'))

# Other: Conda, asdf, Scoop — see https://docs.tilt.dev/install.html

# Verify
tilt version
```

### Prerequisites

```bash
# Required for the Next.js dev image (private @squantumengine/* scope)
export HORIZON_KEY=<your-npm-token>
```

### Usage

```bash
# Start all services with hot-reload
tilt up

# Tilt UI (resource status, logs, manual triggers)
open http://localhost:10350

# Backend health check
curl http://localhost:8000/healthz

# Frontend
open http://localhost:3000

# Stop all services and remove containers
tilt down
```

> **First run is slow** (full image build, same as `docker compose up --build`).
> Subsequent runs reuse Docker layer cache and start in seconds.

---

## Quick start — Docker Compose (production image)

```bash
# Spawn Go service + Postgres + Redis
docker compose up -d postgres redis backend-belajar-lintas-ronald
curl http://localhost:8000/healthz

# Spawn FE service
docker compose up -d frontend-belajar-lintas-ronald
curl http://localhost:3000/

# Teardown
docker compose down -v
```

The compose file uses marker syntax (`# boilerplate:subproject=… axis=… option=…`) so the SDLC plugin pack can scaffold a single-subproject subset for new services. See [`scripts/prune.sh`](scripts/prune.sh) for the prune algorithm and [`docs/superpowers/specs/2026-05-14-consolidate-compose-to-root-design.md`](docs/superpowers/specs/2026-05-14-consolidate-compose-to-root-design.md) for the contract.

## Conventions

- **Commits**: Conventional Commits (`feat(scope): subject`) + `Refs: <TICKET-ID>` footer for `feat`/`fix`. See `.gitmessage`.
- **Branches**: `<type>/<TICKET-ID>-<slug>`. Types: `feature | fix | chore | docs | experiment`.
- **CI path-scoping**: each subproject's CI workflow filters `paths: ['<lang>/**']` so cross-language changes don't trigger unrelated builds.

## Adding a new language

1. Create `<lang>/` at repo root.
2. Add `<lang>/.boilerplate.yaml` (schema v1, see existing manifests for shape).
3. Add `.github/workflows/ci-<lang>.yml` with `paths: ['<lang>/**']`.
4. Add `<lang>/` row to the table in this README.
5. Document inside `<lang>/docs/`.

## Documentation

- Per-subproject product docs:
  - Go: [`backend/docs/`](backend/docs/) — [ARCHITECTURE](backend/docs/ARCHITECTURE.md), [ADAPTERS](backend/docs/ADAPTERS.md), [PRUNING](backend/docs/PRUNING.md), [ERROR_CODES](backend/docs/ERROR_CODES.md), [CONTRIBUTING](backend/docs/CONTRIBUTING.md).
  - Next.js: see [`frontend/README.md`](frontend/README.md) (component rules in [`frontend/CLAUDE.md`](frontend/CLAUDE.md)).
- Monorepo-wide design history (specs + implementation plans): [`docs/superpowers/`](docs/superpowers/).

## Preview Deployment

Ephemeral per-PR environments on GKE, triggered by commenting `/preview` on any PR.
Each environment gets its own namespace (`preview-<repo>-pr-<n>`), subdomain, and isolated infrastructure.

### How it works

1. Comment `/preview` on a PR (OWNER, MEMBER, or COLLABORATOR only)
2. GitHub Actions builds backend + frontend images and pushes to Artifact Registry
3. Helm deploys all resources into a dedicated namespace
4. Migration + seed jobs run automatically after deploy
5. Environment is torn down when the PR is closed

URLs after deploy:

```
App  https://<repo>-pr-<n>.preview.<your-domain>
API  https://<repo>-pr-<n>-api.preview.<your-domain>/healthz
```

### One-time infra setup

| What | Requirement |
|---|---|
| GKE cluster | Traefik ingress controller installed |
| DNS | Wildcard `*.preview.<your-domain>` → cluster ingress LB |
| Artifact Registry | `preview` repository (Docker format) in your GCP project |
| Workload Identity Federation | GitHub Actions → GCP (no long-lived keys) |
| Service account roles | `roles/artifactregistry.writer`, `roles/container.developer` |

### GitHub Secrets

Set in **Settings → Secrets and variables → Actions**:

| Secret | Description |
|---|---|
| `WIF_PROVIDER` | WIF provider resource name |
| `WIF_SERVICE_ACCOUNT` | GCP service account email |
| `GCP_PROJECT_ID` | GCP project ID |
| `AR_REGION` | Artifact Registry region, e.g. `asia-southeast1` |
| `GKE_CLUSTER_NAME` | GKE cluster name |
| `GKE_CLUSTER_LOCATION` | GKE cluster zone/region, e.g. `asia-southeast1-a` |
| `PREVIEW_DB_PASSWORD` | PostgreSQL password for preview instances |
| `HORIZON_KEY` | npm token for `@squantumengine/horizon` (private package) |

### Setup after cloning

All configuration lives in `helm/preview/values.yaml`. Changes needed when using this boilerplate for a new service:

**1. Binary name** — update migration and seed commands to match your built binary:

```yaml
# helm/preview/values.yaml
migration:
  command: ["/your-service-name", "db:migrate"]

seed:
  command: ["/your-service-name", "db:seed"]
```

The binary name comes from `backend/Dockerfile`: `go build -o <name>`.

**2. Domain** — replace `preview.squantumengine.com`:

```yaml
# helm/preview/values.yaml
domain: preview.your-domain.com
```

Also update the domain in `.github/workflows/preview-deploy.yml` (the `NEXT_PUBLIC_BACKEND_URL` build arg and the comment body).

**3. Node pool** — update if your cluster uses different labels:

```yaml
# helm/preview/values.yaml
nodeSelector:
  node_pool: your-node-pool-name

tolerations:
  - key: your-taint-key
    operator: Equal
    value: your-taint-value
    effect: NoSchedule
```

**4. Service ports and health paths** (if different from defaults):

```yaml
# helm/preview/values.yaml
backend:
  port: 8080
  health:
    path: /healthz

frontend:
  port: 3000
  health:
    path: /
```

### Optional services

PostgreSQL and Redis are enabled by default. MongoDB, MySQL, MinIO, and RabbitMQ are disabled by default. Enable any of them in `helm/preview/values.yaml`:

```yaml
mongodb:
  enabled: true
  password: your-password   # required

mysql:
  enabled: true
  password: your-password      # required
  rootPassword: your-root-pass # required

minio:
  enabled: true
  # rootUser and rootPassword default to "minioadmin"

rabbitmq:
  enabled: true
  # username and password default to "rabbit"
```

When a service is enabled, its connection env vars are automatically injected into the backend container with sensible defaults. Override the env var names if your app uses different names:

```yaml
redis:
  extraEnv:
    - name: CACHE_HOST      # rename REDIS_HOST
      value: "redis"
    - name: CACHE_PORT
      value: "6379"
    - name: CACHE_PASSWORD
      value: "previewredispass"
```

### Custom env vars

Inject additional env vars into backend or frontend via `extraEnv` (standard Kubernetes format, supports `valueFrom`):

```yaml
backend:
  extraEnv:
    - name: APP_ENV
      value: "preview"
    - name: JWT_SECRET
      valueFrom:
        secretKeyRef:
          name: my-secret
          key: JWT_SECRET

frontend:
  extraEnv:
    - name: NEXT_PUBLIC_FEATURE_FLAG
      value: "true"
```

### Workflow files

| File | Trigger |
|---|---|
| `.github/workflows/preview-deploy.yml` | `/preview` comment on PR |
| `.github/workflows/preview-cleanup.yml` | PR closed (merged or abandoned) with `preview-deployed` label |

> **Note:** Workflow changes must be merged to `main` before they take effect. The `issue_comment` trigger always runs from `main`.

## License

See `LICENSE` at the repo root.

# Architecture

`backend-belajar-lintas-ronald` uses **Hexagonal Architecture (Ports & Adapters)** with light DDD vocabulary. The goal is to keep business logic isolated from infrastructure so that infrastructure choices (DB, cache, message broker, …) are swap-able.

## Dependency direction

```
        ┌───────────────────────────┐
        │     internal/http/        │  Inbound adapter (Echo)
        └─────────────┬─────────────┘
                      │ uses Service interface
                      ▼
        ┌───────────────────────────┐
        │   internal/domain/<X>/    │  Center
        │   • entities              │  No imports from http/adapter/app
        │   • Service interface     │
        │   • Service impl          │
        │   • Repository interface  │  ← PORT
        │   • Cache interface       │
        └─────────────┬─────────────┘
                      ▲ implements ports
                      │
        ┌─────────────┴─────────────┐
        │   internal/adapter/<Y>/   │  Outbound adapters
        │   • postgres / mysql / …  │
        │   • redis / inmemory / …  │
        │   • otel / datadog / …    │
        └───────────────────────────┘

        ┌───────────────────────────┐
        │       internal/app/       │  Composition root: build adapters,
        │       container.go        │  inject into services, register
        │                           │  routes. The ONLY place that
        │                           │  imports many packages.
        └───────────────────────────┘
```

## Allowed / forbidden imports

Enforced by `golangci-lint depguard` rules in `.golangci.yaml`.

| Allowed | Forbidden |
|---|---|
| `domain` → `pkg/*` | `domain` → `adapter/*`, `http/*`, `app/*` |
| `adapter/<X>` → `domain` (to implement port) | `adapter/<X>` → `adapter/<Y>` (no cross-adapter) |
| `http` → `domain` (to call service) | `http` → `adapter/persistence/*`, `cache/*`, `messaging/*`, … |
| `app` → all | `domain.A` → `domain.B.Repository` (must use `domain.B.Service`) |

Note: `internal/http` may import `internal/adapter/auth/*` — auth middleware bridges the HTTP boundary and is excluded from the `http-no-adapter` rule.

## Cross-domain access

A domain accesses another domain's data **only via that domain's `Service` interface**, never via its `Repository`. This preserves business rules (soft delete, authz, cache, audit, events) that the service enforces.

```go
// internal/domain/order/service.go (hypothetical)
type service struct {
    repo        Repository
    customerSvc customer.Service   // ← cross-domain via interface
}

func (s *service) CreateOrder(ctx context.Context, ...) error {
    customer, err := s.customerSvc.FindByID(ctx, customerID)  // not customerRepo
    // ...
}
```

Wiring in `app/container.go` injects the concrete `customerService` into `orderService`.

### Why not access the repository directly?

Going through the repository directly bypasses business rules carried by the service:

- Soft-delete visibility filters
- Authorization checks against `ctxutil.UserID` / `Roles`
- Cache invalidation on mutation
- Audit logging
- Domain events / outbox

The performance "win" is illusory — Go interface dispatch is essentially free, and the business consequences of bypassing the service are real (security, consistency, observability gaps).

## Composition root

`internal/app/container.go` is the **only** file that imports many adapter packages. It:

1. Reads typed config (`config.Load`)
2. Connects to selected infrastructure (Postgres, Redis, …)
3. Wires adapters into domain services
4. Constructs HTTP handlers and middleware

Marker comments around each axis option allow the SDLC plugin (or `scripts/prune.sh`) to remove blocks for unselected options without manual rewrite.

## Configuration

`internal/config/config.go` declares a single `Config` struct with one field per axis option. Loaded via Viper from `env/env.<APP_ENV>` (defaults to `local`). Each axis option's config block is wrapped in `boilerplate:axis=X option=Y START..END` markers so the plugin can prune unused fields.

## Testing layers

| Layer | Path | Tooling | Coverage target |
|---|---|---|---|
| Unit (domain) | `internal/domain/<x>/*_test.go` | testify + handwritten mocks | ≥ 80% |
| Handler | `internal/http/handler/*_test.go` | httptest + mock service | ≥ 70% |
| Integration (adapter) | `internal/adapter/persistence/<x>/*_test.go` (build tag `integration`) | testcontainers-go | ≥ 60% per adapter |

Run: `make test-unit`, `make test-integration`, `make test-cover`.

## Response envelope

All HTTP responses follow the **lintas team API contract**. The boilerplate ships three envelope shapes:

| Envelope             | Used for                                  | Required fields |
|----------------------|-------------------------------------------|-----------------|
| `SuccessEnvelope`    | Single-resource success + no-content      | `success`, `code`, `message`, `data`, `timestamp` |
| `PaginatedEnvelope`  | List endpoints                            | `SuccessEnvelope` + `pagination` |
| `ErrorEnvelope`      | All error responses                       | `success`, `code` (5-digit numeric), `message`, `timestamp`; optional `metadata` |

The renderer lives at `internal/http/response/envelope.go` and exposes exactly three helpers — `OK`, `Paged`, `Err`. Handlers MUST use them; manual `c.JSON` of envelope-shaped maps defeats the single-source-of-truth guarantee.

### Package boundaries

```
internal/http/handler        → calls response.OK / Paged / Err
internal/http/response       → owns envelope shape; depends on pkg/customerror + pkg/pagination
pkg/customerror              → error primitive (HTTP-agnostic); raisable from any layer
pkg/pagination               → Page struct + paging math (HTTP-agnostic)
```

`pkg/customerror.CustomError` is the cross-layer error currency. Domain, service, and repository layers all raise it. The HTTP layer (`response.Err`) is the only place that translates `CustomError` to wire JSON.

### Numeric error codes

The `ErrorEnvelope.Code` field is a 5-digit integer following the lintas convention `HTTP_PREFIX + 2-digit subcode` (e.g., `40100` = unauthorized variant 00). The boilerplate ships 8 generic codes in `pkg/customerror/error_code.go`; scaffolded services extend the catalog with domain-specific subcodes.

Full convention: see `docs/ERROR_CODES.md`.

### No-content responses

Endpoints that conceptually return "no payload" (e.g., `DELETE`) still return `200 OK` + envelope with `data: {}`. This trades the HTTP-idiomatic `204 No Content` for envelope universality — every client parses the same shape; observability tooling sees the same fields.

`204 No Content` is reserved for non-JSON endpoints (webhooks acks, file streams) — not used by this boilerplate's sample handler.

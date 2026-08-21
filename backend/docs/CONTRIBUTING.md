# Contributing

This document covers two audiences:

1. **Engineers using the boilerplate** — branch/commit/PR conventions for the project derived from this boilerplate.
2. **Contributors to the boilerplate itself** — how to add a new axis or option.

---

## For engineers using the boilerplate

### Branch naming

Format: `<type>/<TICKET-ID>-<slug>`

| Type | Use when |
|---|---|
| `feature` | New feature from a PRD |
| `fix` | Bug fix |
| `chore` | Refactor, dependency update, non-functional change |
| `docs` | Documentation-only change |
| `experiment` | Spike / POC, not intended to merge |

Examples:

- `feature/SDLC-123-payment-flow`
- `fix/SDLC-456-login-redirect`
- `chore/SDLC-789-upgrade-go-1.25`

Rules:

- Slug uses kebab-case.
- Total branch name ≤ 50 chars.
- Direct push to `main` is not allowed.

### Commit messages — Conventional Commits

Format: `<type>(<scope>): <subject>`

```
feat(payment): add stripe checkout integration

Implements card-only flow per PRD section 5.2.
Integration tested against stripe sandbox.

Refs: SDLC-123
```

Valid types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `perf`, `ci`, `build`.

Rules:

- Subject ≤ 72 chars, imperative ("add", not "added").
- Body explains **why**, not **what**.
- Footer `Refs: <TICKET-ID>` is required for `feat`/`fix`.
- Breaking change: append `!` after scope (e.g. `feat(api)!: rename endpoint`).

### Local enforcement

Three pre-commit hooks installed by `make setup`:

1. `commit-msg` — validates Conventional Commits format + `Refs:` footer for feat/fix.
2. `pre-push` — validates branch name.
3. `pre-commit` — runs lint, gitleaks, etc.

### Server enforcement

GitHub Action `.github/workflows/pr-title.yml` validates PR title (matters for squash-merge).

### PR template

`.github/PULL_REQUEST_TEMPLATE.md` — keep summary, link spec, declare type, list test plan.

---

## For boilerplate contributors

### Adding a new option to an existing axis

1. Create the adapter folder under `internal/adapter/<axis>/<new-option>/`.
2. Implement the axis's port interface (e.g., `todo.Cache` for cache axis, the standalone capability interface for messaging/storage/etc.).
3. Add an entry under `axes.<axis>.options.<new-option>` in `.boilerplate.yaml`:
   ```yaml
   <new-option>:
     path: internal/adapter/<axis>/<new-option>
     compose_services: [<services if any>]
     env_vars: [<ENV_VAR_NAMES>]
     go_imports: [<module paths>]
   ```
4. Add `boilerplate:axis=<axis> option=<new-option> START..END` markers in:
   - `internal/app/container.go` (wiring block)
   - `internal/config/config.go` (Config field block + Load() block)
   - `compose.yml` (service block, if any)
   - `env/env.example` (env var block)
   - `Makefile` (target block, if any)
5. `make build && make test` must pass.
6. Update `docs/ADAPTERS.md` with a "When to choose this option" line.
7. Verify markers paired: `bash scripts/verify-prune.sh`.

### Adding a new axis

Heavier — requires:

1. New folder tree `internal/adapter/<new-axis>/<options>/`.
2. New entry under `axes.<new-axis>` in `.boilerplate.yaml`.
3. Wiring code in `internal/app/container.go` (and possibly `server.go`).
4. Config struct field in `internal/config/config.go`.
5. New env vars + markers in `env/env.example`.
6. Compose services + markers in `compose.yml` (if applicable).
7. Documentation entries in `docs/ADAPTERS.md`.
8. If integrated with `todo` reference domain: update `internal/domain/todo/service.go` to take the new port as a constructor dependency, plus update tests + mocks.

### Testing changes

Same as engineers: `make test`, `make test-integration`, `make lint`. The boilerplate has no separate test suite; verification is via the existing test layer.

### Tagging releases

The boilerplate uses SemVer tags (`v1.0.0`, `v1.1.0`, …). Breaking changes (e.g., manifest schema bumps) require a major bump.

`scaffold-be-project` plugin (in the `lintas` plugin pack) clones the latest `vX.Y.Z` tag.

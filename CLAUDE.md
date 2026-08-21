# Boilerplate Monorepo — Working Rules

This file applies when working in the boilerplate-monorepo itself. For
per-stack rules see `backend/CLAUDE.md` or `frontend/CLAUDE.md`. This file
does **not** ship to scaffolded services.

## Repo nature

This is a **template repo**, not a runnable service. The lintas SDLC
plugin pack (`sdlc-be-go`, `sdlc-fe-nextjs`, future `sdlc-be-python` /
`sdlc-be-java`) clones a subfolder at a tagged version and prunes axis
options the team did not pick.

→ Overview: `README.md`
→ Per-subproject overview: `backend/README.md`, `frontend/README.md`

## Rule 1 — Axis-parity guardrail

Every change preserves compile + runtime validity for **all** axis
options, not just the default. The boilerplate's value is the matrix
of working combinations.

Verification matrix:
- `cd golang && go build ./...` — clean exit
- `cd nextjs && pnpm tsc --noEmit` — clean across at least three FE
  auth swap states: `jwt-refresh`, `opaque-session`, `none`
- `cd golang && go build ./...` clean for `OBSERVABILITY_BACKEND` ∈
  {otel, datadog, noop} (unpruned tree covers all three simultaneously)
- `cd nextjs && pnpm tsc --noEmit` clean in the unpruned tree and in
  each observability prune state: `aws-rum`, `sentry`, `otel`, `noop`

→ Detail: `docs/superpowers/specs/2026-05-14-fe-auth-multi-adapter-design.md`
→ Manifests: `backend/.boilerplate.yaml`, `frontend/.boilerplate.yaml`

## Rule 2 — Marker discipline

Wiring files use marker comments:

```
// boilerplate:axis=X option=Y START
... block content ...
// boilerplate:axis=X option=Y END
```

When importing axis-scoped content from a non-scoped file, marker-wrap
**both** the import declaration AND every usage. Reference pattern:
`frontend/src/mocks/handlers.ts`.

```typescript
// boilerplate:axis=sample-app option=todo-app START
import { todosHandlers } from './todos';
// boilerplate:axis=sample-app option=todo-app END

export const handlers = [
  // boilerplate:axis=sample-app option=todo-app START
  ...todosHandlers,
  // boilerplate:axis=sample-app option=todo-app END
];
```

→ Pruner algorithm: `backend/docs/PRUNING.md`
→ Subproject marker dimension: `docs/superpowers/specs/2026-05-14-consolidate-compose-to-root-design.md`

## Rule 3 — Manifest fields beyond `path`

Axis option entries in `.boilerplate.yaml` support more than a single
`path`:

```yaml
todo-app:
  path: src/app/(protected)/todos          # primary folder to delete
  extra_paths:
    - src/mocks/todos                      # additional folders to delete
  codemods:
    - file: src/static/route.ts
      replace: defaultBaseUrlPage value '/todos' -> '/' when todo-app unselected
  env_vars: []
```

Use `extra_paths` when an option spans multiple folders. Use
`codemods` when an option requires literal-value replacement rather
than wholesale folder deletion. Pruner consumes both.

→ Schema: `backend/docs/PRUNING.md`
→ Live example: `frontend/.boilerplate.yaml` — `sample-app.options.todo-app`

## Rule 4 — Branch & commit convention

- Branch: `<type>/<TICKET-ID>-<slug>`. Types: `feature`, `fix`,
  `chore`, `docs`, `experiment`. Slug uses kebab-case.
- Commit: Conventional Commits with `Refs: <TICKET-ID>` footer
  required for `feat` and `fix` types.
- Never use `--no-verify` to bypass hooks. Investigate hook failures.

→ Hook detail: `backend/docs/CONTRIBUTING.md`
→ Commit template: `.gitmessage`

## Rule 5 — Documentation tier model

Three tiers of CLAUDE.md:

| Tier | File | Ships to scaffolded service? | Audience |
|---|---|---|---|
| 1 | `CLAUDE.md` (this file) | No | Monorepo developers only |
| 2 | `backend/CLAUDE.md` | Yes | Go service developers |
| 2 | `python/CLAUDE.md` | Yes | Python service developers |
| 3 | `frontend/CLAUDE.md` | Yes | Next.js service developers |
| 3 | `qa-test/CLAUDE.md` | Yes | QA automation developers |

Per-stack files (tiers 2-3) MUST be self-contained — they cannot
reference this root file because it does not exist in the scaffolded
service.

→ Design history: `docs/superpowers/specs/2026-05-14-3tier-claude-md-design.md`

## Rule 6 — Language preference

- Chat output: Bahasa Indonesia (per lintas plugin user preference).
- Generated artifact files (PRD.md, tech-spec.md, api-contract.yaml,
  test-plan.md): English per repo convention.
- Commit messages, branch names, code comments, documentation files
  (this file, all files in `docs/`): English.

## Quick reference

| Topic | Source |
|---|---|
| Subproject overview | `backend/README.md`, `frontend/README.md` |
| Per-stack rules | `backend/CLAUDE.md`, `frontend/CLAUDE.md` |
| Pruning algorithm | `backend/docs/PRUNING.md` |
| Compose marker dimension | `docs/superpowers/specs/2026-05-14-consolidate-compose-to-root-design.md` |
| Design history | `docs/superpowers/specs/`, `docs/superpowers/plans/` |
| Commit template | `.gitmessage` |

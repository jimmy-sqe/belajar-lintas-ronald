# Contributing

This document covers two audiences:

1. **Engineers using the boilerplate** — branch/commit/PR conventions for
   the Next.js project derived from this boilerplate.
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

- `feature/SDLC-123-todo-bulk-edit`
- `fix/SDLC-456-login-redirect`
- `chore/SDLC-789-upgrade-react-19`

Rules:

- Slug uses kebab-case.
- Total branch name ≤ 50 chars.
- Direct push to `main` is not allowed.

### Commit messages — Conventional Commits

Format: `<type>(<scope>): <subject>`

```
feat(todos): add bulk-edit modal

Implements the multi-select + batch-update pattern per PRD section 4.2.
Reuses existing TodoRow component for selection rendering.

Refs: SDLC-123
```

Valid types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`,
`chore`, `perf`, `ci`, `build`.

Rules:

- Subject ≤ 72 chars, imperative ("add", not "added").
- Body explains **why**, not **what**.
- Footer `Refs: <TICKET-ID>` is required for `feat`/`fix`.
- Breaking change: append `!` after scope (e.g., `feat(auth)!: rename token endpoint`).

### Local enforcement

Husky hook installed by `pnpm install` (via the `prepare` script):

1. `pre-commit` — auto-formats staged files via `prettier --write` and re-stages them.

Note: there is no `commit-msg` hook in `.husky/`. Conventional Commits format is NOT enforced locally by hooks; CI runs the format check on PR title (`.github/workflows/pr-title.yml`).

### Required local commands before pushing

```bash
pnpm build              # ensures Horizon package builds clean (CI gates on this)
pnpm lint:fix           # auto-fix lint
pnpm prettier --write . # format
# pnpm test            # (no test script wired yet — FE test infra is a planned follow-up)
```

CI will fail and block the pipeline if `pnpm build` fails or lint /
type checks regress.

### PR template

`.github/PULL_REQUEST_TEMPLATE.md` — keep summary, link spec/PRD,
declare type, list test plan, paste screenshot for UI changes.

---

## For boilerplate contributors

### Adding a new option to an existing axis

For **folder-prune** axes (auth, sample-app, mocking):

1. Create the option folder under the axis's conventional path.
2. Satisfy the axis's interface contract:
   - `auth` — implement `AuthAdapter` from
     `src/services/auth/types.ts`.
   - `sample-app` — self-contained folder under `(auth)/` or
     `(protected)/` with its own layout + styles.
   - `mocking` — MSW handlers under `src/mocks/<option>/`.
3. Add an entry under `axes.<axis>.options.<new-option>` in
   `.boilerplate.yaml`:
   ```yaml
   <new-option>:
     path: src/...
     extra_paths: [...]   # if multi-folder
     codemods: [...]      # if literal-replacement needed
     env_vars: [...]
     description: ...
   ```
4. Add `boilerplate:axis=<axis> option=<new-option> START..END`
   markers in any wiring file that imports the new option (e.g.,
   `src/mocks/handlers.ts`, `src/services/auth/index.ts`,
   `src/app/layout.tsx`).
5. Verify `pnpm tsc --noEmit` clean.
6. Update `docs/AXES.md` with a "When to choose" row.

For **preset** axes (state, forms, validation, ui-library, etc.):

1. Add an entry under `axes.<axis>.options.<new-option>` with
   `npm_deps` listed.
2. Add codemod hints for `layout.tsx` provider replacement if
   relevant.
3. Update `docs/AXES.md`.

### Adding a new axis

Heavier — requires:

1. Decide pruning style (folder-prune vs preset).
2. New folder tree (for folder-prune) or new manifest entry only
   (for preset).
3. New entry under `axes.<new-axis>` in `.boilerplate.yaml`.
4. Wiring code (provider, middleware, hook) plus markers in the
   relevant wiring file from `wiring_files:` list.
5. Documentation entries in `docs/AXES.md` and
   `docs/ARCHITECTURE.md` (if the axis affects the layered model).
6. If integrated with the `todo-app` sample-app: update
   `src/app/(protected)/todos/` and tests.

### Testing changes

Same as engineers: `pnpm build && pnpm lint:fix`. There is no
separate boilerplate test suite; FE test infrastructure
(Jest/RTL/Playwright) is a planned follow-up.

### Tagging releases

The boilerplate uses SemVer tags (`v1.0.0`, `v1.1.0`, ...). Breaking
changes (e.g., manifest schema bumps) require a major bump. The
`sdlc-fe-nextjs:scaffold-fe-project` plugin clones the latest
`vX.Y.Z` tag.

# Axes Catalog

The boilerplate ships **11 modular axes** for Next.js. The canonical
schema is `.boilerplate.yaml`; this document explains when to choose
each option.

## Multiplicity rules

- **`exactly_one`** — must select exactly one option. A `none` option
  is provided as the explicit "off" choice.
- **`one_or_more`** — select at least `min` options.

## Style of pruning

| Style | How the plugin handles it | Examples |
|---|---|---|
| **Folder-prune** | Delete the option's `path` (and `extra_paths`); pruner removes marker-wrapped blocks in wiring files. | `auth`, `sample-app`, `mocking` |
| **Preset** | No folder deletion. Plugin manipulates `package.json` deps and applies codemods to `layout.tsx` providers. Boilerplate ships with the default option already wired. | `state`, `forms`, `validation`, `ui-library`, `api-client`, `styling`, `observability`, `rendering` |

This is unlike the Go boilerplate, where every axis is folder-prune.

## Defaults shipped

```yaml
rendering: csr
auth: jwt-refresh
api-client: kubb
state: zustand
forms: react-hook-form
validation: zod
styling: [tailwind]
ui-library: horizon
mocking: msw
observability: none
```

## rendering — `exactly_one` (preset)

| Option | When to choose |
|---|---|
| `csr` | Default. Encrypted localStorage; axios interceptor handles refresh. |
| `ssr` | Server-side cookies (httpOnly encrypted); middleware rotates refresh tokens. Choose for SEO-critical pages or strict CSP environments. |

## auth — `exactly_one` (folder-prune)

Each option lives at `src/services/auth/<option>/`. App code imports
only `@/services/auth`.

| Option | When to choose |
|---|---|
| `jwt-refresh` | Default. Access + refresh tokens; client-side expiry decode; auto-refresh queue. |
| `opaque-session` | Server-side opaque token; bearer header; 401 → re-login (no client-side decode). |
| `oauth2-oidc` | Stub. Not implemented in v1.0.0. |
| `none` | Public app. Stubs throw "auth disabled". |

## sample-app — `one_or_more` (folder-prune, min 0)

Opinionated demo features. Each option's folder is self-contained
(includes its own AppShell + scoped CSS) and removable wholesale.

| Option | Path | Requires | When to choose |
|---|---|---|---|
| `auth-flow` | `src/app/(auth)` | `auth != none` | Login, change-password, forgot-password (forgot is jwt-refresh only). |
| `todo-app` | `src/app/(protected)/todos` | `auth != none` | TodoApp CRUD pages (list, detail, create/edit/delete modals) — reference sample-app implementation. |

### `extra_paths` and `codemods` example

`sample-app.todo-app` declares additional concerns beyond its primary
`path`:

```yaml
todo-app:
  path: src/app/(protected)/todos
  extra_paths:
    - src/mocks/todos                # also delete this folder
  codemods:
    - file: src/static/route.ts
      replace: defaultBaseUrlPage value '/todos' -> '/' when todo-app unselected
```

Use these fields when an option spans multiple folders or requires
literal-value replacement instead of folder deletion.

## api-client — `exactly_one` (preset)

| Option | When to choose |
|---|---|
| `kubb` | Default. Spec-first via OpenAPI → types + Zod schemas + React Query hooks. Output in `src/openapi/` (never hand-edited). |
| `orval` | Stub. Not implemented in v1.0.0. |
| `hand-written` | Stub. Not implemented in v1.0.0. |

## state — `exactly_one` (preset)

| Option | When to choose |
|---|---|
| `zustand` | Default. Small store under `src/store/`. One slice per domain. |
| `jotai` | Stub. Atom-based. Not implemented in v1.0.0. |
| `redux-toolkit` | Stub. Heavier; for complex state needs. Not implemented in v1.0.0. |
| `none` | React `useState` only. Plugin removes `src/store/`. |

## forms — `exactly_one` (preset)

| Option | When to choose |
|---|---|
| `react-hook-form` | Default. Used in `(auth)/` pages and TodoApp modals. |
| `none` | Plain controlled inputs. |

## validation — `exactly_one` (preset)

| Option | When to choose |
|---|---|
| `zod` | Default. Used by Kubb output and the React Hook Form resolver. |
| `yup` | Stub. Not implemented in v1.0.0. |
| `none` | No runtime validation. |

## styling — `one_or_more` (preset, min 1)

| Option | When to choose |
|---|---|
| `tailwind` | Default. Layout and spacing utility classes. |
| `css-modules` | Built-in to Next.js; no extra config. Opt-in via `*.module.css`. |
| `scss` | Already in `package.json`. Opt-in via `*.scss` files. |

## ui-library — `exactly_one` (preset)

| Option | When to choose |
|---|---|
| `horizon` | Default. Internal SQE component library. `nextjs/CLAUDE.md` Rule 1 mandates it as the sole UI library. |
| `shadcn` | Stub. Open-source component generator. Not implemented in v1.0.0. |
| `mui` | Stub. Material UI. Not implemented in v1.0.0. |
| `custom` | Roll your own. |

## mocking — `exactly_one` (folder-prune)

| Option | When to choose |
|---|---|
| `msw` | Default. Mock Service Worker for dev + test. Handlers in `src/mocks/`. |
| `none` | Hit the real backend always. Plugin removes `src/mocks/` and its imports. |

## observability — `exactly_one` (preset)

| Option | When to choose |
|---|---|
| `aws-rum` | AWS CloudWatch RUM. Dependency present, no wiring yet. |
| `sentry` | Stub. Not implemented in v1.0.0. |
| `otel` | Stub. Not implemented in v1.0.0. |
| `none` | Default. No client observability. |

## Adding a new option to an existing axis

1. Create the option folder (for folder-prune axes) under the axis's
   conventional path (e.g., `src/services/auth/<new-option>/` or
   `src/app/(protected)/<new-sample-app>/`).
2. Satisfy the axis's interface contract (e.g., `AuthAdapter` for
   auth; see `docs/ARCHITECTURE.md#auth-adapters`).
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
4. Add `boilerplate:axis=<axis> option=<new-option> START..END` markers
   in any wiring file that references the new option.
5. Run `pnpm tsc --noEmit` to confirm compile parity.
6. Update this `docs/AXES.md` with a "When to choose" row.

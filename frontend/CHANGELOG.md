# Changelog

## [Unreleased]

### Added — TodoApp minimum runnable

- **`src/app/(protected)/todos/`**: list + detail pages with create/edit/delete modals (Horizon + RHF + Zod). Shared `todoFormSchema` for create + edit.
- **`src/services/todos/`**: stable internal hook names (`useTodos`, `useTodo`, `useCreateTodo`, `useUpdateTodo`, `useDeleteTodo`) wrapping kubb-generated hooks. Page code consumes only these — kubb output rename doesn't propagate.
- **`src/mocks/todos/handlers.ts`**: MSW handlers with lintas envelope responses + in-memory list mutation across POST/PUT/DELETE.

### Changed (BREAKING — TodoApp)

- **`auth` axis** (jwt-refresh adapter): response types adjusted to consume lintas envelope shape `{success, code, message, data: {...}, timestamp}`. `signIn` / `refresh` parse `res.data.data.*`. `auth.ts` now uses `customFetch` default export with axios-style config (replaces previously broken `import { fetcher }`).
- **`SignInRequest`** field renamed `username` → `email` (PRD: email + password auth). `useLoginPage.ts` + `login/page.tsx` updated; post-login redirect targets `/todos` (was `/dashboard`).
- **`sample-app` axis**: option `dashboard` renamed → `todo-app`, path `src/app/(protected)/todos`.
- **`src/openapi/openapi.json`** rewritten to mirror BE `openapi.yaml` (auth + todos paths). Kubb regenerate produces new `src/openapi/{types,zod}/` plus flat operation hooks at `src/openapi/` root.
- **`pnpm-lock.yaml`** resynced with `package.json` (drift: `@ant-design/icons`/`nextjs-registry` removed, `@svgr/webpack` added).

### Removed

- `(protected)/dashboard/` page folder (placeholder superseded by todos).
- `src/mocks/items/` (empty placeholder).
- `JwtChangePasswordResponse` + `JwtResetPasswordResponse` types and corresponding MSW handlers (out of TodoApp MVP scope).

Plugin SDLC `sdlc-fe-nextjs` in `lintas` repo must be updated to handle the renamed sample-app option `todo-app` (separate follow-up).

### Changed (BREAKING — auth + sample-app)

- **Auth restructure.** `src/services/auth/` is now a folder per auth model: `jwt-refresh/`, `opaque-session/`, `oauth2-oidc/`, `none/`. App code consumes only `@/services/auth` (the index re-export). Direct deep imports from `@/services/auth/<option>/...` are forbidden.
- **Auth types hand-written.** Each option's `types.ts` defines its response shape directly. `src/openapi/openapi.json` no longer contains `/v1/auth/*` paths. Kubb-generated `HttpCommonResponseWrapper` is replaced by hand-written `HttpErrorResponse` in `src/services/types.ts`.
- **`auth` axis renamed**: option `jwt-cookie-or-localstorage` → `jwt-refresh`. New option `opaque-session` added (for backends using opaque session tokens with no refresh, e.g. COSMOS).
- **`sample-app` axis added** (`one_or_more`, `min=0`): controls inclusion of `(auth)/` pages and `(protected)/dashboard/`. Production scaffolds can opt out.
- **`forgot-password` page** gated by `.boilerplate-prune.yaml` directive — kept only when `auth=jwt-refresh` AND `sample-app` contains `auth-flow`.
- **Deleted**: `src/services/{auth.ts,session.ts}`, `src/services/{client,server}/auth.ts`, `src/services/server/cookies.ts`, `src/services/fetcher/refresh-queue.ts`, `src/static/types/session.d.ts`. All logic now lives under `src/services/auth/jwt-refresh/`.
- **MSW reorganised**: per-model handlers under `src/mocks/auth/<option>/`; `src/mocks/handlers.ts` aggregates with namespaced imports + marker blocks.

### Fixed

- `not-found.tsx` no longer references a missing illustration; uses text-only 404 panel.

### Migration notes

- Plugin SDLC `sdlc-fe-nextjs` must be updated in the `lintas` repo to parse the new axis options and the sub-folder `.boilerplate-prune.yaml` directive. Until that update lands, scaffolded projects may carry unused auth-model folders; users can manually `rm -rf src/services/auth/<unused>/`.
- Existing scaffolded projects are unaffected; their generated state is independent of this boilerplate refactor.

## [v1.0.0] — 2026-05-13

### Added

- **Next.js 15 App Router** boilerplate on React 19 + TypeScript 5.
- **Dual SSR / CSR auth**: `NEXT_PUBLIC_RENDERING_MODE` switches storage strategy.
- **Generated API client** via Kubb: types/, zod/, React Query hooks emitted from `src/openapi/openapi.json`.
- **MSW mocking** toggled by `NEXT_PUBLIC_API_MOCKING`.
- **Horizon UI** (`@squantumengine/horizon`) as the sole UI library (see `CLAUDE.md`).
- **Zustand** for state, **react-hook-form** + **Zod** for forms.
- **Tailwind CSS** for layout/spacing.
- **Husky** pre-commit hook running lint + prettier on staged FE files.
- **`.boilerplate.yaml` manifest** declaring 10 modular axes for the SDLC plugin.
- **`docs/{ARCHITECTURE,ADAPTERS,CONTRIBUTING}.md`** documentation.
- **Generic OpenAPI 3.0.3 stub** at `src/openapi/openapi.json` (auth + items CRUD, parity with the Go boilerplate's reference domain).
- **Root CI workflow** `.github/workflows/ci-nextjs.yml` filtered to `paths: ['nextjs/**']`.

### Changed (from initial commit)

- `package.json` name: `smart-shipping-fe` → `frontend-belajar-lintas-ronald`, version reset to `1.0.0`.
- Typo `src/config/environtment.ts` → `environment.ts` (+ 12 import paths).
- Theme: `theme-hz-mining` → `theme-hz-blue` (generic Horizon theme).
- Default redirect: `/admin` → `/dashboard`; route group `(protected)/admin/` → `(protected)/dashboard/`.
- Metadata: `Create Next App` → `frontend-belajar-lintas-ronald`.

### Removed

- `@ant-design/icons` and `@ant-design/nextjs-registry` from `package.json` (contradicted `CLAUDE.md`).
- SQE-specific generated API content: Partner-token endpoints, OTP-based change-password.
- `nextjs/.git/` (nested repo) and `nextjs/.github/` (GitHub only reads root `.github/`); content migrated to repo root.

### Notes

- Engineers should replace `src/openapi/openapi.json` with their backend's spec and run `pnpm generate:schema`.
- The boilerplate ships with default selection (csr + horizon + zustand + react-hook-form + zod + tailwind + msw). Axes are documented in `.boilerplate.yaml`.

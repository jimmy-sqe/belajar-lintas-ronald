# Architecture

`frontend-belajar-lintas-ronald` uses the **Next.js 15 App Router** with route
groups for access-level separation and **folder-per-option** for
axis-scoped concerns. The goal is to keep features self-contained
so the SDLC plugin can prune unselected options by deleting whole
folders.

## Folder layout

```
nextjs/src/
├── app/                          App Router tree
│   ├── layout.tsx                Provider stack (Theme, Toaster,
│   │                              ReactQuery, MSW, Session)
│   ├── design-tokens.css         Horizon color/spacing tokens (GLOBAL)
│   ├── globals.css               Tailwind base + global resets
│   ├── (auth)/                   PUBLIC route group
│   │   └── login/
│   │       └── _styles/login.css  Scoped CSS for auth pages
│   └── (protected)/              AUTHENTICATED route group
│       ├── layout.tsx            Minimal wrapper (no AppShell here)
│       └── todos/                Sample-app feature, self-contained
│           ├── layout.tsx        Sidebar + Topbar + AppShell live HERE
│           ├── page.tsx          /todos list
│           ├── [id]/page.tsx     /todos/:id detail
│           ├── _components/      Modals, rows, empty states
│           ├── _hooks/           Page-level hooks
│           └── _styles/todoapp.css  Scoped CSS
├── components/                   Cross-feature presentational components
├── hooks/                        Cross-feature hooks
├── services/                     Feature-agnostic clients
│   ├── auth/                     Folder-per-option auth adapter
│   │   ├── index.ts              Re-exports the selected option
│   │   ├── jwt-refresh/          Default option
│   │   ├── opaque-session/       Server-side opaque session
│   │   ├── oauth2-oidc/          Stub (not implemented in v1.0.0)
│   │   └── none/                 Stub for public apps
│   └── fetcher/                  Axios + auth refresh integration
├── store/                        Zustand stores (one slice per domain)
├── openapi/                      Kubb-generated types + hooks (DO NOT edit)
├── mocks/                        MSW handlers
│   ├── handlers.ts               Aggregator, uses marker-import pattern
│   ├── auth/<option>/            Per-option auth mocks
│   └── todos/                    sample-app=todo-app mocks
├── utils/                        Utility helpers (zod schemas, etc.)
├── config/                       Env-driven feature flags
└── static/                       Static config (route table, brand)
```

## Dependency direction

```
   ┌──────────────────────────────────┐
   │     src/app/<route-group>/       │  Pages, layouts (Server + Client)
   └────────────┬─────────────────────┘
                │ imports
                ▼
   ┌──────────────────────────────────┐
   │  src/components, src/hooks       │  Presentational + behavior
   └────────────┬─────────────────────┘
                │ imports
                ▼
   ┌──────────────────────────────────┐
   │  src/services, src/store         │  Side-effect boundary
   │  src/openapi (Kubb hooks)        │
   └──────────────────────────────────┘
```

Pages and layouts MAY import from any of `components`, `hooks`,
`services`, `store`, `openapi`. They SHOULD NOT import deep into
another route group.

## Route groups

Next.js route groups `(name)/` partition the app by access level
without affecting URLs:

- `(auth)/login` → `/login` (public)
- `(protected)/todos` → `/todos` (authenticated)

Each route group owns its own `layout.tsx`. The minimal
`(protected)/layout.tsx` exists only to enforce session middleware;
it does NOT host AppShell. Each sample-app feature (e.g., todos)
owns its own AppShell inside `(protected)/<feature>/layout.tsx`. This
keeps prune-by-folder clean: removing `(protected)/todos/` deletes
its sidebar, topbar, and styling in one operation without leaving
orphan chrome.

## Wiring files (axis pruning operates here)

The SDLC plugin removes marker-wrapped blocks in these files when an
axis option is unselected:

| File | What gets marker-wrapped |
|---|---|
| `src/app/layout.tsx` | Provider blocks (per axis option) |
| `src/middleware.ts` | SSR refresh-token guard (per rendering option) |
| `src/config/environment.ts` | Env-driven feature flags |
| `package.json` | `npm_deps` per axis |
| `.env-example` | Env vars per axis |
| `src/mocks/handlers.ts` | Per-feature mock imports (marker-import pattern) |

Note: `src/mocks/handlers.ts` is not in the manifest's `wiring_files:` list but uses the same marker-import pattern (see Rule 5 of `CLAUDE.md`). The manifest may be extended to include it in a future version.

## Provider stack

`src/app/layout.tsx` wraps the entire tree in this order
(outer-to-inner):

```
MSWProvider
  └── ToasterProvider (Horizon)
        └── ReactQueryProvider
              └── SessionProvider (auth)
                    └── {children}
```

Theming is applied via a `theme-hz-blue` CSS class on `<body>`, NOT a JSX `<ThemeProvider>` wrapper. The actual Horizon theme token sheet is imported globally.

Modifying any provider order requires updating this section + the
boilerplate axis manifest.

## Auth adapter contract

All `src/services/auth/<option>/` folders satisfy the same
`AuthAdapter` interface declared in `src/services/auth/types.ts`.
The exact interface as of this boilerplate version:

```typescript
export interface AuthAdapter {
  signIn(request: SignInRequest): Promise<Session>;
  signOut(): Promise<void>;
  getSession(): Promise<Session | null>;
  /** Only for models that support refresh (jwt-refresh). */
  refresh?: (session: Session) => Promise<Session>;
  changePassword?: (request: ChangePasswordRequest) => Promise<void>;
  /** Only for models that support self-service reset (jwt-refresh). */
  resetPassword?: (request: ResetPasswordRequest) => Promise<void>;
}
```

Where `SignInRequest = { email: string; password: string }`,
`ChangePasswordRequest = { currentPassword: string; newPassword: string }`,
and `ResetPasswordRequest = { email: string }` (all declared in
`src/services/auth/types.ts`).

`src/services/auth/index.ts` re-exports the adapter instance for the
selected option:

```typescript
import adapter from './jwt-refresh';
export default adapter;
export * from './jwt-refresh/types';
export * from './types';
```

The plugin scaffold operates here by replacing the literal string
`'./jwt-refresh'` with the chosen option's folder name. No marker
comments needed for this file because the substitution is
literal-based.

App code consumes ONLY `@/services/auth` — never the option subfolder
directly.

## Server vs Client components

- Default: Server Component (no `'use client'` directive).
- Add `'use client'` only when the file uses hooks, browser APIs,
  event handlers, or client-only libraries (Zustand stores, RHF,
  React Query mutations).
- Provider components in `layout.tsx` are Client by necessity.
- Server actions (`'use server'`) are NOT used in this boilerplate
  (auth flows via REST endpoints over fetch).

## Testing layers

| Layer | Path | Tooling |
|---|---|---|
| Unit | `<file>.test.tsx` next to subject | Jest + React Testing Library |
| Integration (component + MSW) | same path, uses MSW handlers | Jest + MSW |
| E2E | `e2e/<feature>.spec.ts` | Playwright (when configured) |

Test name encodes FR ID: `describe('CreateTodoModal — FR005', ...)`.

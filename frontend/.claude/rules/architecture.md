# Architecture Rules

Detail for the architecture directives summarized in `CLAUDE.md`:
route groups, the auth adapter, the marker-import pattern, and the
reference sample-app feature.

## App Router route groups

- `(auth)/` — public auth pages (login, change-password,
  forgot-password). No layout chrome.
- `(protected)/` — authenticated pages. May host multiple sample-app
  features.
- Each route group owns its own `layout.tsx`.

Sample-app features (e.g., todos) are **self-contained** inside their
route group folder. Sidebar, Topbar, AppShell, and per-feature CSS
live INSIDE `(protected)/<feature>/`, NOT in parent
`(protected)/layout.tsx`. This keeps prune-by-folder clean — removing
a sample-app option deletes its entire folder without leaving orphan
chrome.

→ Detail: `docs/ARCHITECTURE.md#route-groups`

## Folder-per-option auth adapter

`src/services/auth/<option>/` (jwt-refresh, opaque-session,
oauth2-oidc, none). Each folder satisfies the same `AuthAdapter`
interface.

- App code imports from `@/services/auth` ONLY. NEVER deep-import
  `@/services/auth/jwt-refresh/...` from outside the adapter folder.
- The exported `AuthAdapter` instance is re-exported by
  `src/services/auth/index.ts` for the selected option.
- Stub options (`oauth2-oidc`, `none`) may throw "not implemented" at
  runtime but MUST compile and satisfy the interface contract.

→ Adapter contract: `docs/ARCHITECTURE.md#auth-adapters`

## Marker-import pattern

When you import from an axis-scoped folder INTO a file that itself is
not axis-scoped (e.g., `src/mocks/handlers.ts` imports from
`src/mocks/todos/`), marker-wrap BOTH the import declaration AND
every usage:

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

The pruner strips both blocks; remaining code compiles clean.

Reference: `src/mocks/handlers.ts`.

## Reference domain

`src/app/(protected)/todos/` is the end-to-end reference sample-app
feature. Use as pattern when adding a new sample-app option. Pages:
list (`page.tsx`), detail (`[id]/page.tsx`), modals
(`_components/`), shared layout (`layout.tsx`), scoped CSS
(`_styles/todoapp.css`).

# Quality Rules — FR annotation, testing, build verification

Detail for the quality-gate directives summarized in `CLAUDE.md`.

## FR annotation pattern

Components, hooks, or pages that implement a functional requirement
carry a `// FR-<ID>: <one-line description>` comment above the
declaration:

```typescript
// FR-005: TodoApp create modal with form validation
export function CreateTodoModal({ ... }: Props) { ... }
```

Lintas QA skills parse this annotation to compute coverage.

## Testing

- Unit: Jest + React Testing Library. File: `<file>.test.tsx` next to
  the subject.
- E2E: Playwright when configured.
- Backend mocking: MSW handlers in `src/mocks/`.
- Test name encodes FR: `describe('CreateTodoModal — FR005', ...)`.

→ Run: `pnpm test`, `pnpm e2e` (when configured).

## Build verification

After any UI change, ALWAYS run in this order:

```bash
pnpm build       # ensures Horizon package builds clean (CI gates on this)
pnpm lint:fix    # auto-fix lint issues
pnpm prettier --write .   # format
```

CI will fail and block the pipeline if `pnpm build` fails.

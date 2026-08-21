# Next.js Service — Working Rules

This file applies when working in this Next.js service (boilerplate or
scaffolded). Loaded automatically by Claude Code.

This is a **thin index**: each rule below states the core directive,
then points (`→`) to a detail doc under `.claude/rules/` for examples,
code, and elaboration. Read the detail doc when a rule is relevant —
don't guess at the parts kept out of this file.

## Stack at a glance

Next.js 15 (App Router) · React 19 · TypeScript 5 · Tailwind ·
Horizon UI · Zustand · React Hook Form + Zod · Kubb (OpenAPI codegen)
· MSW.

→ Full architecture: `docs/ARCHITECTURE.md`
→ Axis catalog (11 axes): `docs/AXES.md`
→ Manifest: `.boilerplate.yaml`

## UI & components — `.claude/rules/ui.md`

- **Horizon is the sole component library.** NEVER hand-roll a
  component that exists in Horizon, NEVER use raw `<button>`/`<input>`/
  `<select>`/`<table>`/`<dialog>` when Horizon covers it, and NEVER
  install another UI library — except the two carve-outs below.
- **ALWAYS** import from `@squantumengine/horizon`; NEVER use the
  `hz-` prefix in our own `className=`.
- **Read `.claude/docs/horizon-components.md` before building UI** —
  do not guess prop names. Verify icon names against the Horizon icon
  catalog.
- **Theming/providers** are wired in `src/app/layout.tsx` via a
  `theme-hz-<color>` body class (NOT a JSX `<ThemeProvider>`) —
  PRESERVE on edits.
- **Charts → Apache ECharts; maps → Leaflet** (the only sanctioned
  non-Horizon UI libs). Both are client-only: `'use client'` +
  `next/dynamic({ ssr: false })`.

→ Full directives, quick-lookup component table, theming details,
charts/maps setup, and the Figma-to-code workflow: `.claude/rules/ui.md`

## Architecture — `.claude/rules/architecture.md`

- **Route groups:** `(auth)/` (public, no chrome) and `(protected)/`
  (authenticated). Each owns its own `layout.tsx`. Sample-app features
  are self-contained inside `(protected)/<feature>/` so prune-by-folder
  stays clean.
- **Auth adapter:** folder-per-option under `src/services/auth/<option>/`.
  App code imports from `@/services/auth` ONLY — NEVER deep-import an
  adapter folder.
- **Marker-import pattern:** when importing axis-scoped content into a
  non-scoped file, marker-wrap BOTH the import AND every usage.

→ Route groups, adapter contract, marker example, and the `todos/`
reference feature: `.claude/rules/architecture.md`

## Data, state & forms — `.claude/rules/data.md`

- **State:** Zustand for client-owned state only (`src/store/`). Server
  state goes through Kubb-generated React Query hooks — do NOT
  duplicate server data into Zustand.
- **API client:** Kubb-generated into `src/openapi/`. NEVER hand-edit;
  regenerate via `pnpm kubb`. Consume only via generated hooks.
- **Forms:** React Hook Form + Zod; prefer Horizon `FormTextField`.

→ Detail: `.claude/rules/data.md`

## Quality — `.claude/rules/quality.md`

- **FR annotation:** mark code that implements a requirement with
  `// FR-<ID>: <description>` (QA skills parse it for coverage).
- **Testing:** Jest + RTL (`<file>.test.tsx`), Playwright E2E when
  configured, MSW for backend mocks. Test names encode the FR.
- **Build verification:** after any UI change, run `pnpm build` →
  `pnpm lint:fix` → `pnpm prettier --write .`. CI gates on `pnpm build`.

→ Detail: `.claude/rules/quality.md`

## Workflow — `.claude/rules/workflow.md`

- **Conventional Commits**; `feat`/`fix` REQUIRE a `Refs: <TICKET-ID>`
  footer. Branch `<type>/<TICKET-ID>-<slug>`. NEVER `--no-verify`.

→ Detail: `.claude/rules/workflow.md`

## Quick reference

| Topic | Source |
|---|---|
| UI / Horizon / charts / maps / Figma | `.claude/rules/ui.md` |
| Route groups / auth adapter / markers | `.claude/rules/architecture.md` |
| State / API client / forms | `.claude/rules/data.md` |
| FR annotation / testing / build | `.claude/rules/quality.md` |
| Commit / PR / branch | `.claude/rules/workflow.md` |
| Horizon component API | `.claude/docs/horizon-components.md` |
| Horizon icon catalog | https://sqehorizon.squantumengine.com/?path=/docs/foundation-icon--overview |
| Architecture + route groups | `docs/ARCHITECTURE.md` |
| Axis catalog | `docs/AXES.md` |
| Commit/PR/hook detail | `docs/CONTRIBUTING.md` |
| Manifest | `.boilerplate.yaml` |
| Pruning algorithm (monorepo) | `../golang/docs/PRUNING.md` |

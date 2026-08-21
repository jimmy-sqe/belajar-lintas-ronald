# Pruning Contract

The boilerplate ships every option of every axis. Real projects use **only** the chosen
combination. Two consumers prune the boilerplate after scaffolding:

1. **Claude Code SDLC plugin** (`sdlc-fe-nextjs`) — reads `.boilerplate.yaml` + PRD/TSD,
   programmatically selects options, removes everything else. Primary consumer.
2. **Manual pruning** — a developer follows this contract by hand. Same logic,
   manual execution.

## Pruning styles

Next.js axes use two distinct pruning strategies, unlike the Go and Python boilerplates
where every axis is folder-prune.

| Style | How the plugin handles it | Axes |
|---|---|---|
| **Folder-prune** | Deletes `option.path` (and `extra_paths`); strips marker blocks in wiring files. | `auth`, `sample-app`, `mocking` |
| **Preset** | No folder deletion. Plugin removes `npm_deps` from `package.json`, applies `codemods` to swap barrel imports (e.g., `services/auth/index.ts`). Boilerplate ships the default option already wired. | `rendering`, `state`, `forms`, `validation`, `ui-library`, `api-client`, `styling`, `observability` |

## Manifest: `.boilerplate.yaml`

The canonical machine-readable contract. Top-level keys:

- `version` — schema version (currently `1`)
- `name`, `description` — boilerplate identity
- `subproject`, `folder` — scaffold placement hints
- `axes` — map of axis name to axis definition
- `core_files` — globs that always exist regardless of selection
- `wiring_files` — files that contain marker blocks and must be edited (not deleted) during prune

Per axis:

- `multiplicity` — `exactly_one` or `one_or_more`
- `min` (only for `one_or_more`) — minimum selection count
- `options` — map of option name to option definition

Per option:

- `path` — folder to delete if not selected (folder-prune axes only; `null` for preset axes)
- `extra_paths` — additional files/folders to delete with the option
- `npm_deps` — npm package names the option adds to `package.json`. Removed for unselected
  options by Stage G; packages shared with a selected option are preserved.
- `codemods` — literal find/replace rules applied to wiring files (e.g., swapping the auth
  adapter import in `src/services/auth/index.ts`)
- `env_vars` — env vars referenced in `.env-example`

## Marker syntax

Comment-bracketed blocks, regex-friendly, line-based:

| File type | Marker |
|---|---|
| TypeScript / TSX / JavaScript (`*.ts`, `*.tsx`, `*.js`) | `// boilerplate:axis=X option=Y START` … `// boilerplate:axis=X option=Y END` |
| YAML (`.yml`, `.yaml`) | `# boilerplate:axis=X option=Y START` … `# boilerplate:axis=X option=Y END` |
| JSON (`package.json`, etc.) | Markers are **not** used inside JSON. npm_deps are handled via Stage G instead. |
| `.env*` | `# boilerplate:axis=X option=Y START` … `# boilerplate:axis=X option=Y END` |

Rules:

- Markers must be **paired** — every START has an END. `pnpm tsc --noEmit` failing after
  prune is the signal that a marker block was malformed or left dangling.
- Markers may **NOT nest**.
- A marker block contains content for **one** option only — never mix.
- When importing an axis-scoped module into a non-scoped file, marker-wrap **both** the
  import declaration AND every usage (see `src/mocks/handlers.ts` as the reference pattern).

## Prune algorithm

```
INPUT: selection_map (per .boilerplate.yaml)

Stage A — validate
1. Validate selections against multiplicity. Reject if invalid.
2. Compute UNSELECTED options = all_options \ selection_map.

Stage B — folder-prune axes (auth, sample-app, mocking)
3. For each UNSELECTED option:
   a. Delete `option.path` recursively (if non-null).
   b. Delete each entry in `extra_paths`.
   c. From every `wiring_file`, remove all lines from
      `// boilerplate:axis=X option=Y START` to
      `// boilerplate:axis=X option=Y END` (inclusive).
4. For each SELECTED option:
   a. From every `wiring_file`, remove ONLY the marker comment lines themselves
      (the START and END lines) — keep the code between them.

Stage C — compose services
5. Compute used compose_services = union across SELECTED options.
   For every compose service NOT in used_services, remove its block from
   `compose.yml` using its marker.

Stage D — env vars
6. Compute used env_vars; prune `.env-example` accordingly.

Stage F — codemods (preset axes)
7. For each option, apply its `codemods[]` rules. Used to repoint folder-per-option
   barrel files (e.g., `src/services/auth/index.ts` → selected auth option folder,
   `src/services/observability/index.ts` → selected observability option folder).

Stage G — npm_deps (preset axes)
8. Remove each UNSELECTED option's `npm_deps` from `package.json`
   (dependencies / devDependencies / peerDependencies). Preserve any dep also
   listed by a SELECTED option.

Stage H — verify
9. Run:
     pnpm install
     pnpm tsc --noEmit
     pnpm build
10. If all pass, commit: "chore: prune unused adapters per <selection>".
```

## codemods

Some axes (all preset axes, plus the folder-prune `auth` axis for `src/services/auth/index.ts`)
require literal value replacement rather than wholesale block deletion.

Live examples in `.boilerplate.yaml`:

- `auth.*` — swaps the import path in `src/services/auth/index.ts` from `'./jwt-refresh'` to
  the chosen option folder name
- `observability.*` — swaps the import in `src/services/observability/index.ts` and
  `src/services/observability/server.ts`
- `rendering.ssr` — enables the `src/middleware.ts` SSR refresh guard

Codemod format in `.boilerplate.yaml`:

```yaml
codemods:
  - file: src/services/auth/index.ts
    replace: "'./jwt-refresh'" -> "'./opaque-session'" when auth=opaque-session selected
```

## Wiring files

Files that contain markers and must be edited (not deleted) during prune. Reproduced from
`.boilerplate.yaml`:

- `src/app/layout.tsx` — provider blocks (per axis option)
- `src/middleware.ts` — SSR refresh-token guard (per rendering option)
- `src/config/environment.ts` — env-driven feature flags
- `src/mocks/handlers.ts` — per-feature mock imports (marker-import pattern)
- `.env-example` — env var blocks per axis

Note: `package.json` is handled by Stage G (npm_deps), not by marker blocks. JSON does not
support line comments, so `npm_deps` is the mechanism for that file.

## Verification

After pruning, verify:

1. `pnpm tsc --noEmit` — all TypeScript compiles. This catches dangling marker blocks and
   broken imports in one command.
2. `pnpm build` — Next.js production build succeeds. CI gates on this.
3. `pnpm lint` — no lint errors introduced.
4. Markers are paired — every remaining `// boilerplate:axis=X option=Y START` has a
   matching `END`.

`pnpm tsc --noEmit` is the analogue of `go build ./...` / `make compileall` on other stacks.
Any failure means the prune was incomplete or a marker block was malformed.

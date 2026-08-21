# Pruning Contract

The boilerplate ships every option of every axis. Real projects use **only** the chosen combination. Two consumers prune the boilerplate after scaffolding:

1. **Claude Code SDLC plugin** — reads `.boilerplate.yaml` + PRD/TSD, programmatically selects options, removes everything else. Primary consumer.
2. **`scripts/prune.sh`** — interactive shell pruner for human users. Same logic, manual prompts.

## Manifest: `.boilerplate.yaml`

The canonical machine-readable contract. Top-level keys:

- `version` — schema version (currently `1`)
- `name`, `description` — boilerplate identity
- `axes` — map of axis name → axis definition
- `core_files` — globs that always exist regardless of selection
- `wiring_files` — files that contain marker blocks and must be edited (not deleted) during prune

Per axis:

- `multiplicity` — `exactly_one` or `one_or_more`
- `min` (only for `one_or_more`) — minimum selection count
- `options` — map of option name → option definition

Per option:

- `path` — folder to delete if not selected (or `null` for option-without-folder, e.g. `auth.none`)
- `migrations_path` — optional, deleted with the option
- `compose_services` — service names referenced in `compose.yml` (deduplicated across axes)
- `env_vars` — env vars referenced in `env/env.example`
- `go_imports` — Go module paths the option uses (informational; `go mod tidy` is the source of truth)
- `codemods` — literal find/replace rules (`file`, `pattern`, `replacement`, `when_selected` / `when_unselected`). Applied by Stage F.
- `npm_deps` — npm package names the option adds to `package.json` (nextjs subproject). Removed for unselected options by Stage G; packages shared with a selected option are preserved.
- `makefile_targets` (api-doc only) — Makefile target names guarded by markers
- `files` (container axis only) — file paths the option owns

## Marker syntax

Comment-bracketed blocks, regex-friendly, line-based:

| File type | Marker |
|---|---|
| Go (`.go`) | `// boilerplate:axis=X option=Y START` … `// boilerplate:axis=X option=Y END` |
| YAML (`.yml`, `.yaml`) | `# boilerplate:axis=X option=Y START` … `# boilerplate:axis=X option=Y END` |
| Makefile / `.env*` | `# boilerplate:axis=X option=Y START` … `# boilerplate:axis=X option=Y END` |

Rules:

- Markers must be **paired** — every START has an END. Verified by `scripts/verify-prune.sh`.
- Markers may **NOT nest**.
- A marker block contains content for **one** option only — never mix.

## Prune algorithm

```
INPUT: selection_map (per .boilerplate.yaml)

1. Validate selections against multiplicity. Reject if invalid.
2. Compute UNSELECTED options = all_options \ selection_map.
3. For each UNSELECTED option:
   a. Delete `option.path` recursively (if non-null).
   b. Delete `migrations_path` if applicable.
   c. From every `wiring_file`, remove all lines from
      `boilerplate:axis=X option=Y START` to `boilerplate:axis=X option=Y END` (inclusive).
4. For each SELECTED option:
   a. From every `wiring_file`, remove ONLY the marker comment lines themselves
      (the START and END lines) — keep the wiring code between them.
5. Compute used compose_services = union across SELECTED options.
   For every compose service NOT in used_services, remove its block from
   `compose.yml` using its marker.
6. Compute used env_vars; prune `env/env.example` accordingly.
6b. Stage F (codemods): for each option, apply its `codemods[]` rules —
    `when_selected` rules fire when the option IS selected,
    `when_unselected` when it is NOT. Used to repoint folder-per-option
    barrels (e.g. src/services/auth/index.ts,
    src/services/observability/index.ts + server.ts).
6c. Stage G (npm_deps): for the nextjs subproject, remove each
    UNSELECTED option's `npm_deps` from package.json
    (dependencies/devDependencies/peerDependencies), preserving any dep
    also listed by a SELECTED option. No-op when no package.json exists.
7. Run:
     go mod tidy
     go build ./...
     golangci-lint run
     go test -short ./...
8. If all pass, commit: "chore: prune unused adapters per <selection>".
```

## Wiring files

Files that contain markers and must be edited (not deleted) during prune:

- `internal/app/container.go`
- `internal/app/server.go`
- `internal/config/config.go`
- `cmd/migrate.go`
- `compose.yml`
- `env/env.example`
- `Makefile`

## Sanity check: `scripts/verify-prune.sh`

After pruning, run:

```bash
bash scripts/verify-prune.sh
```

Checks:

1. Every `START` in wiring files has a matching `END`.
2. No orphan markers in non-wiring files.
3. `go build ./...` passes.
4. `go mod tidy` is a no-op (no leftover unused deps).

Exit non-zero if any check fails.

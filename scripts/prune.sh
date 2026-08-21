#!/usr/bin/env bash
# scripts/prune.sh — Compose pruner for the monorepo.
#
# Reads a root compose.yml that uses dimension-aware markers:
#   # boilerplate:subproject=<list> [axis=<NAME> option=<NAME>] START
#   <content>
#   # boilerplate:subproject=<list> [axis=<NAME> option=<NAME>] END
#
# Filters the compose for a single target subproject + that subproject's
# axis option selections, then emits the pruned compose. Marker lines are
# stripped from kept blocks.
#
# Usage:
#   scripts/prune.sh --subproject=<name> [--selections=<axis=opt,axis=opt,...>] \
#                    [--compose=<path>] [--output=<path>]
#
# Defaults:
#   --compose=compose.yml at repo root
#   --output=- (stdout)
#
# Requirements: bash 4+, yq, grep, sed.

set -euo pipefail

# ---------- Bash version guard ----------
if (( BASH_VERSINFO[0] < 4 )); then
    echo "ERROR: bash 4+ required (found ${BASH_VERSION}); on macOS install via 'brew install bash'" >&2
    exit 1
fi

# ---------- Arg parsing ----------
# New v1 flags (in addition to legacy --subproject/--selections/--output)
MANIFEST_FILE=""   # --manifest=<subproject|path>: read subproject+selections from manifest itself
DRY_RUN=0

# Resolved in main_full_pipeline from manifest .folder key.
# SELECTION_SUBPROJECT = name used in boilerplate:subproject= markers
# SELECTION_FOLDER     = directory on disk for filesystem operations
SELECTION_FOLDER=""

# Legacy CSV-mode flags (preserved for backward-compat)
SUBPROJECT=""
SELECTIONS=""
COMPOSE_FILE="compose.yml"
OUTPUT="-"

usage() {
    # Help text goes to stdout when exit code is 0 (--help), stderr otherwise.
    local code="${1:-2}"
    local fd=2
    [[ "$code" == "0" ]] && fd=1
    cat >&"$fd" <<USAGE
Usage:
  $0 --manifest=<subproject> [--dry-run]
      Full pipeline: compose + adapters + wiring + env + Makefile + codemods + verify.
      Reads subproject, folder, and selections from the subproject's own
      .boilerplate.yaml (requires a 'selections:' block added by /design-tech-spec).

  $0 --subproject=<name> [--selections=<axis=opt,axis=opt,...>]
        [--compose=<path>] [--output=<path>]
      Legacy compose-only mode (preserved for backward-compat).

v1 flags:
  --manifest=<subproject>    Subproject name (e.g., golang, nextjs, python) whose
                             .boilerplate.yaml contains a 'selections:' block.
                             Accepts a full path to the manifest as well.
  --dry-run                  Preview plan, do not modify filesystem

Legacy flags:
  --subproject=<name>        Target subproject (e.g., golang, nextjs, python)
  --selections=<list>        Comma-separated axis=option pairs
                             (e.g., persistence=postgres,cache=redis)
  --compose=<path>           Path to source compose (default: compose.yml)
  --output=<path>            Output path or "-" for stdout (default: -)
USAGE
    exit "$code"
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --manifest=*)        MANIFEST_FILE="${1#*=}" ;;
        --dry-run)           DRY_RUN=1 ;;
        --subproject=*)      SUBPROJECT="${1#*=}" ;;
        --selections=*)      SELECTIONS="${1#*=}" ;;
        --compose=*)         COMPOSE_FILE="${1#*=}" ;;
        --output=*)          OUTPUT="${1#*=}" ;;
        --help|-h)           usage 0 ;;
        *) echo "Unknown argument: $1" >&2; usage 2 ;;
    esac
    shift
done

# Dispatch flag: legacy CSV mode if --selections=<csv> given without --manifest.
# Not used in any control flow yet — wired in later tasks.
if [[ -n "$SELECTIONS" && -z "$MANIFEST_FILE" ]]; then
    LEGACY_MODE=1
else
    LEGACY_MODE=0
fi

# Normalize --manifest value: accept "golang" (subproject name) or a full path.
# Expand to "<name>/.boilerplate.yaml" when no path separator present.
if [[ -n "$MANIFEST_FILE" && "$MANIFEST_FILE" != */* && "$MANIFEST_FILE" != *.yaml ]]; then
    MANIFEST_FILE="${MANIFEST_FILE}/.boilerplate.yaml"
fi

# Mode-specific basic validation.
# Full mode: --manifest given.
# Legacy mode: --subproject required + compose file must exist.
MANIFEST=""
if [[ -n "$MANIFEST_FILE" ]]; then
    # Full mode: defer subproject/manifest resolution to main_full_pipeline.
    :
else
    # Legacy mode validation
    [[ -z "$SUBPROJECT" ]] && { echo "ERROR: --subproject is required" >&2; usage 2; }
    [[ -f "$COMPOSE_FILE" ]] || { echo "ERROR: compose file not found: $COMPOSE_FILE" >&2; exit 3; }

    MANIFEST="${SUBPROJECT}/.boilerplate.yaml"
    if [[ ! -f "$MANIFEST" ]]; then
        # Fixture mode: look in scripts/test/fixtures/ for the manifest
        ALT_MANIFEST="scripts/test/fixtures/${SUBPROJECT}.manifest.yml"
        if [[ -f "$ALT_MANIFEST" ]]; then
            MANIFEST="$ALT_MANIFEST"
        else
            echo "ERROR: manifest not found at $MANIFEST" >&2
            exit 4
        fi
    fi
fi

# ---------- Tool dependency check (after all argument validation) ----------
command -v yq >/dev/null 2>&1 || { echo "ERROR: yq is required" >&2; exit 5; }

# jq is required for stage_npm_deps (nextjs subproject). Soft check: warn
# if missing; fail later at stage_npm_deps invocation if actually needed.
if ! command -v jq >/dev/null 2>&1; then
    echo "[WARN] jq not found — stage_npm_deps will fail if any subproject has package.json" >&2
fi

# SEL is the global axis→option map. Populated by parse_selection_csv
# (legacy mode) or parse_selection_file (full mode) — both defined below.
declare -A SEL=()

# ====================================================================
# SECTION 2: Manifest helpers (yq wrappers for .boilerplate.yaml)
# ====================================================================
manifest_has_axis() {
    yq -e ".axes.\"$2\"" "${1:-$MANIFEST}" >/dev/null 2>&1
}

manifest_axes() {
    yq -r '.axes | keys | .[]' "${1:-$MANIFEST}"
}

manifest_options() {
    yq -r ".axes.\"$2\".options | keys | .[]" "${1:-$MANIFEST}"
}

manifest_axis_multiplicity() {
    yq -r ".axes.\"$2\".multiplicity // \"exactly_one\"" "${1:-$MANIFEST}"
}

manifest_axis_min() {
    yq -r ".axes.\"$2\".min // 1" "${1:-$MANIFEST}"
}

manifest_option_path() {
    local out
    out=$(yq -r ".axes.\"$2\".options.\"$3\".path // \"\"" "${1:-$MANIFEST}")
    [[ "$out" == "null" ]] && out=""
    printf '%s' "$out"
}

manifest_option_migrations_path() {
    local out
    out=$(yq -r ".axes.\"$2\".options.\"$3\".migrations_path // \"\"" "${1:-$MANIFEST}")
    [[ "$out" == "null" ]] && out=""
    printf '%s' "$out"
}

manifest_option_seeds_path() {
    local out
    out=$(yq -r ".axes.\"$2\".options.\"$3\".seeds_path // \"\"" "${1:-$MANIFEST}")
    [[ "$out" == "null" ]] && out=""
    printf '%s' "$out"
}

manifest_option_extra_paths() {
    yq -r ".axes.\"$2\".options.\"$3\".extra_paths // [] | .[]" "${1:-$MANIFEST}"
}

# subproject_in_list <target> <comma_list>  → 0 if in list, 1 if not
subproject_in_list() {
    local target="$1"
    local list="$2"
    IFS=',' read -ra arr <<< "$list"
    for x in "${arr[@]}"; do
        [[ "$x" == "$target" ]] && return 0
    done
    return 1
}

# ====================================================================
# SECTION 3: Selection parser (CSV legacy + YAML new)
# ====================================================================

# parse_selection_csv <csv>
# Populates global SEL associative array. Format: axis=opt,axis=opt,...
parse_selection_csv() {
    local csv="$1"
    declare -gA SEL=()
    [[ -z "$csv" ]] && return 0
    IFS=',' read -ra _pairs <<< "$csv"
    for _p in "${_pairs[@]}"; do
        local _k="${_p%%=*}"
        local _v="${_p#*=}"
        if [[ -z "$_k" || -z "$_v" || "$_k" == "$_p" ]]; then
            echo "ERROR: invalid selection entry: '$_p' (expected axis=option)" >&2
            return 6
        fi
        SEL["$_k"]="$_v"
    done
}

# parse_selection_file <yaml_file>
# Populates global SEL (axis → option) and SELECTION_SUBPROJECT.
parse_selection_file() {
    local file="$1"
    [[ -f "$file" ]] || { echo "ERROR: selections file not found: $file" >&2; return 1; }
    yq -e '.' "$file" >/dev/null 2>&1 || { echo "ERROR: invalid YAML: $file" >&2; return 1; }

    SELECTION_SUBPROJECT=$(yq -r '.subproject // ""' "$file")
    [[ -n "$SELECTION_SUBPROJECT" ]] || {
        echo "ERROR: missing required field 'subproject' in $file" >&2
        return 1
    }

    declare -gA SEL=()
    while IFS=$'\t' read -r axis option; do
        if [[ -n "$axis" ]]; then
            SEL["$axis"]="$option"
        fi
    done < <(yq -r '.selections | to_entries[] | [.key, .value] | @tsv' "$file")
}

# ====================================================================
# SECTION 5: Validation
# ====================================================================

# validate_selection <manifest_file>
# Returns 0 if SEL is consistent with manifest. Exits with code 3 on error.
# Side effect: implicit 'none' assignment for one_or_more axes with min:0
# that have no selection.
validate_selection() {
    local manifest="$1"
    [[ -f "$manifest" ]] || { echo "ERROR: manifest not found: $manifest" >&2; return 1; }

    local errors=0
    local axes_list
    axes_list=$(manifest_axes "$manifest")

    # Check 5: every key in SEL must be a declared axis
    for axis in "${!SEL[@]}"; do
        if ! grep -qFx "$axis" <<< "$axes_list"; then
            echo "ERROR: unknown axis '$axis' (not declared in $manifest)" >&2
            errors=$((errors+1))
        fi
    done

    # Check 6+7: each value valid; multiplicity satisfied
    for axis in "${!SEL[@]}"; do
        local opt="${SEL[$axis]}"
        local mult min options_list
        mult=$(manifest_axis_multiplicity "$manifest" "$axis")
        min=$(manifest_axis_min "$manifest" "$axis")
        options_list=$(manifest_options "$manifest" "$axis")

        # 'none' is the sentinel for "no selection" ONLY when (a) axis is
        # one_or_more with min:0 AND (b) manifest does NOT declare a literal
        # 'none' option. If manifest declares 'none' as a real option, treat
        # it as a normal value.
        if [[ "$opt" == "none" ]] && ! grep -qFx "none" <<< "$options_list"; then
            if [[ "$mult" == "one_or_more" && "$min" == "0" ]]; then
                continue
            else
                echo "ERROR: option 'none' invalid for axis '$axis' (multiplicity=$mult, min=$min)" >&2
                errors=$((errors+1))
                continue
            fi
        fi

        if ! grep -qFx "$opt" <<< "$options_list"; then
            echo "ERROR: unknown option '$opt' for axis '$axis'. Valid: $(echo "$options_list" | tr '\n' ' ')" >&2
            errors=$((errors+1))
        fi
    done

    # Check 8: every required axis has a selection (or implicit 'none')
    while IFS= read -r axis; do
        if [[ -z "${SEL[$axis]:-}" ]]; then
            local mult min
            mult=$(manifest_axis_multiplicity "$manifest" "$axis")
            min=$(manifest_axis_min "$manifest" "$axis")
            if [[ "$mult" == "one_or_more" && "$min" == "0" ]]; then
                SEL["$axis"]="none"
            else
                echo "ERROR: axis '$axis' has no selection (required by manifest)" >&2
                errors=$((errors+1))
            fi
        fi
    done <<< "$axes_list"

    if (( errors > 0 )); then
        return 3
    fi
}

# ---------- Generalized marker filter ----------
#
# filter_markers <input_file> <comment_prefix> <subproject_filter> [manifest_file]
#
# Drop-in generalization of the inline compose-pruning logic below. Preserves:
#   * Stack-based nesting (multiple START frames may nest; line is emitted only
#     when every frame on the stack votes KEEP).
#   * START/END fingerprint validation (mismatched END triggers exit 7).
#   * Rule E2 leniency: when a manifest is provided and the axis is NOT declared
#     by that manifest, the block is kept by default.
#
# Args:
#   $1 input_file        — file to filter
#   $2 comment_prefix    — "#" for YAML/env/Makefile/compose, "//" for Go
#   $3 subproject_filter — current subproject name; empty string disables the
#                          subproject dimension (subproject= portion of marker
#                          becomes optional and is not used to gate KEEP)
#   $4 manifest_file     — optional; when set, Rule E2 leniency is enabled by
#                          probing axes via yq. Empty/omitted = strict mode
#                          (unknown axis with no SEL match → drop).
#
# Globals consumed:
#   SEL[]  — associative array, axis → option, set by caller
#
# Exit codes:
#   0 on success
#   7 on START/END mismatch, orphan END, or unterminated block
#
# Output: filtered content to stdout.
filter_markers() {
    local input_file="$1"
    local comment_prefix="$2"
    local subproject_filter="$3"
    local manifest_file="${4:-}"

    # Comment prefix (e.g. "#" or "//") is interpolated literally. Two regex
    # variants handle the two filter modes:
    #   * with_subproj: subproject= required, axis/option optional
    #   * no_subproj:   subproject= optional, axis/option required
    local re_with_subproj="^[[:space:]]*${comment_prefix}[[:space:]]+boilerplate:subproject=([A-Za-z0-9_,-]+)([[:space:]]+axis=([A-Za-z0-9_-]+)[[:space:]]+option=([A-Za-z0-9_-]+))?[[:space:]]+(START|END)$"
    local re_no_subproj="^[[:space:]]*${comment_prefix}[[:space:]]+boilerplate:(subproject=([A-Za-z0-9_,-]+)[[:space:]]+)?axis=([A-Za-z0-9_-]+)[[:space:]]+option=([A-Za-z0-9_-]+)[[:space:]]+(START|END)$"

    local -a stack_subproj=() stack_axis=() stack_option=() stack_keep=() stack_marker=()
    local line line_no=0
    local subproj_list axis option marker_kind
    local local_keep depth top i ok

    # emit only if every frame votes KEEP
    _filter_emit_if_all_keep() {
        local j
        for (( j=0; j<${#stack_keep[@]}; j++ )); do
            (( stack_keep[j] == 0 )) && return 0
        done
        printf '%s\n' "$1"
    }

    while IFS= read -r line || [[ -n "$line" ]]; do
        line_no=$((line_no + 1))

        subproj_list=""
        axis=""
        option=""
        marker_kind=""

        if [[ -n "$subproject_filter" && "$line" =~ $re_with_subproj ]]; then
            subproj_list="${BASH_REMATCH[1]}"
            axis="${BASH_REMATCH[3]}"
            option="${BASH_REMATCH[4]}"
            marker_kind="${BASH_REMATCH[5]}"
        elif [[ -z "$subproject_filter" && "$line" =~ $re_no_subproj ]]; then
            subproj_list="${BASH_REMATCH[2]}"   # may be empty
            axis="${BASH_REMATCH[3]}"
            option="${BASH_REMATCH[4]}"
            marker_kind="${BASH_REMATCH[5]}"
        fi

        if [[ -n "$marker_kind" ]]; then
            if [[ "$marker_kind" == "START" ]]; then
                local_keep=0
                # Subproject dimension: only enforced if filter is set.
                ok=1
                if [[ -n "$subproject_filter" ]]; then
                    if ! subproject_in_list "$subproject_filter" "$subproj_list"; then
                        ok=0
                        # Foreign subproject block: pass content through unchanged.
                        # Markers are still consumed; only the content is emitted.
                        # This ensures a two-pass prune (backend then frontend) does
                        # not accidentally delete the other subproject's services.
                        local_keep=1
                    fi
                fi
                if (( ok == 1 )); then
                    if [[ -z "$axis" ]]; then
                        local_keep=1
                    else
                        # Rule E2 leniency requires a manifest. Without one, fall
                        # through to strict SEL matching.
                        if [[ -n "$manifest_file" ]] \
                                && ! yq -e ".axes.\"$axis\"" "$manifest_file" >/dev/null 2>&1; then
                            # axis not declared → keep
                            local_keep=1
                        elif [[ "${SEL[$axis]:-}" == "$option" ]]; then
                            local_keep=1
                        fi
                    fi
                fi
                stack_subproj+=("$subproj_list")
                stack_axis+=("$axis")
                stack_option+=("$option")
                stack_keep+=("$local_keep")
                stack_marker+=("$line")
                continue
            else
                # END marker
                depth=${#stack_keep[@]}
                if (( depth == 0 )); then
                    echo "ERROR: END marker without START at line $line_no" >&2
                    unset -f _filter_emit_if_all_keep
                    return 7
                fi
                top=$(( depth - 1 ))
                if [[ "$subproj_list" != "${stack_subproj[$top]}" \
                      || "$axis" != "${stack_axis[$top]}" \
                      || "$option" != "${stack_option[$top]}" ]]; then
                    echo "ERROR: END at line $line_no does not match START '${stack_marker[$top]}'" >&2
                    unset -f _filter_emit_if_all_keep
                    return 7
                fi
                unset 'stack_subproj[top]' 'stack_axis[top]' 'stack_option[top]' 'stack_keep[top]' 'stack_marker[top]'
                continue
            fi
        fi

        # Non-marker line
        if (( ${#stack_keep[@]} == 0 )); then
            printf '%s\n' "$line"
        else
            _filter_emit_if_all_keep "$line"
        fi
    done < "$input_file"

    if (( ${#stack_keep[@]} != 0 )); then
        top=$(( ${#stack_keep[@]} - 1 ))
        echo "ERROR: unterminated marker block (opened by '${stack_marker[$top]}')" >&2
        unset -f _filter_emit_if_all_keep
        return 7
    fi

    unset -f _filter_emit_if_all_keep
    return 0
}

# ====================================================================
# SECTION 6: Stage executors
# ====================================================================

# stage_compose <compose_file> <subproject> <output> <manifest_file>
#
# Filters root compose.yml via filter_markers and writes to the requested
# output (stdout when output == "-", else atomic write via mktemp+mv).
# Returns 7 if filter_markers reports a marker error.
stage_compose() {
    local compose_file="$1"
    local subproject="$2"
    local output="$3"
    local manifest_file="$4"

    [[ -f "$compose_file" ]] || {
        echo "ERROR: compose file not found: $compose_file" >&2
        return 3
    }

    if [[ "$output" == "-" ]]; then
        filter_markers "$compose_file" "#" "$subproject" "$manifest_file"
        return $?
    fi

    local tmp rc
    tmp=$(mktemp)
    if filter_markers "$compose_file" "#" "$subproject" "$manifest_file" > "$tmp"; then
        mv "$tmp" "$output"
        return 0
    else
        rc=$?
        rm -f "$tmp"
        return "$rc"
    fi
}

# stage_adapters <manifest_file> <subproject_root>
#
# For each UNSELECTED option across all axes in the manifest, removes:
#   - option.path (e.g., internal/adapter/persistence/mysql)
#   - option.migrations_path (e.g., db/migrations/mysql)
#   - option.seeds_path (e.g., db/seeds/mysql)
#   - every entry in option.extra_paths[]
# Missing-path entries are logged as warnings (manifest may declare paths
# that don't apply to the current state, e.g. .down.sql files when policy
# is forward-only).
stage_adapters() {
    local manifest="$1"
    local root="$2"
    local deleted=0 missing=0

    while IFS= read -r axis; do
        local selected="${SEL[$axis]:-}"
        while IFS= read -r option; do
            [[ "$option" == "$selected" ]] && continue

            local path mig sds
            path=$(manifest_option_path "$manifest" "$axis" "$option")
            mig=$(manifest_option_migrations_path "$manifest" "$axis" "$option")
            sds=$(manifest_option_seeds_path "$manifest" "$axis" "$option")

            local targets=()
            [[ -n "$path" ]] && targets+=("$path")
            [[ -n "$mig" ]]  && targets+=("$mig")
            [[ -n "$sds" ]]  && targets+=("$sds")
            while IFS= read -r extra; do
                [[ -n "$extra" ]] && targets+=("$extra")
            done < <(manifest_option_extra_paths "$manifest" "$axis" "$option")

            for t in "${targets[@]}"; do
                local full="$root/$t"
                if [[ -e "$full" ]]; then
                    rm -rf "$full"
                    deleted=$((deleted + 1))
                else
                    echo "[WARN]  manifest declares missing path: $full" >&2
                    missing=$((missing + 1))
                fi
            done
        done < <(manifest_options "$manifest" "$axis")
    done < <(manifest_axes "$manifest")

    echo "[INFO]  Stage B (adapters): deleted $deleted paths, $missing missing"
}

# manifest_wiring_files <manifest_file> → newline-separated wiring file paths
manifest_wiring_files() {
    yq -r '.wiring_files // [] | .[]' "${1:-$MANIFEST}"
}

# _comment_prefix_for <filename> → comment prefix, or empty if not markerable
_comment_prefix_for() {
    case "$1" in
        *.go|*.ts|*.tsx|*.js|*.jsx|*.mjs|*.cjs)  printf '%s' "//" ;;
        *.yaml|*.yml|Makefile|*.env*|.env*)       printf '%s' "#" ;;
        *.json)  printf '%s' "" ;;   # JSON: no comments; skip in stage_wiring
        *)       printf '%s' "#" ;;
    esac
}

# stage_wiring <manifest_file> <subproject_root>
#
# Iterates manifest.wiring_files, runs filter_markers on each. Comment
# prefix selected by file extension. Env file and Makefile are also
# wiring files (in the manifest list) so this single function handles
# stages C, D, and E of the spec.
stage_wiring() {
    local manifest="$1"
    local root="$2"
    local edited=0 missing=0

    while IFS= read -r rel; do
        [[ -z "$rel" ]] && continue
        local full="$root/$rel"
        if [[ ! -f "$full" ]]; then
            echo "[WARN]  wiring file declared but missing: $full" >&2
            missing=$((missing + 1))
            continue
        fi
        local prefix
        prefix=$(_comment_prefix_for "$rel")
        if [[ -z "$prefix" ]]; then
            echo "[INFO]  wiring: skipping $rel (not markerable; handled by another stage)"
            continue
        fi
        local tmp
        tmp=$(mktemp)
        if filter_markers "$full" "$prefix" "" "$manifest" > "$tmp"; then
            if ! cmp -s "$full" "$tmp"; then
                mv "$tmp" "$full"
                edited=$((edited + 1))
            else
                rm -f "$tmp"
            fi
        else
            local rc=$?
            rm -f "$tmp"
            echo "[ERROR] filter_markers failed on $full (rc=$rc)" >&2
            return "$rc"
        fi
    done < <(manifest_wiring_files "$manifest")

    echo "[INFO]  Stage C/D/E (wiring/env/Makefile): edited $edited files, $missing missing"
}

# stage_codemods <manifest_file> <subproject_root>
#
# Iterates every option's codemods[] array. Each rule has:
#   file:           path relative to subproject_root
#   pattern:        sed pattern (s|<pattern>|<replacement>|g)
#   replacement:    sed replacement
#   when_selected:  apply only if this option IS selected (bool, default false)
#   when_unselected:apply only if this option is NOT selected (bool, default false)
# If neither when_* is true, the rule never fires (explicit-better-than-implicit).
stage_codemods() {
    local manifest="$1"
    local root="$2"
    local applied=0

    while IFS= read -r axis; do
        local selected="${SEL[$axis]:-}"
        while IFS= read -r option; do
            local n
            n=$(yq -r ".axes.\"$axis\".options.\"$option\".codemods // [] | length" "$manifest")
            local i=0
            while [[ $i -lt $n ]]; do
                local file pattern replacement when_sel when_unsel
                file=$(yq -r ".axes.\"$axis\".options.\"$option\".codemods[$i].file" "$manifest")
                pattern=$(yq -r ".axes.\"$axis\".options.\"$option\".codemods[$i].pattern" "$manifest")
                replacement=$(yq -r ".axes.\"$axis\".options.\"$option\".codemods[$i].replacement" "$manifest")
                when_sel=$(yq -r ".axes.\"$axis\".options.\"$option\".codemods[$i].when_selected // false" "$manifest")
                when_unsel=$(yq -r ".axes.\"$axis\".options.\"$option\".codemods[$i].when_unselected // false" "$manifest")

                local should=0
                if [[ "$when_sel" == "true" && "$selected" == "$option" ]]; then
                    should=1
                elif [[ "$when_unsel" == "true" && "$selected" != "$option" ]]; then
                    should=1
                fi

                if [[ $should -eq 1 ]]; then
                    local target="$root/$file"
                    if [[ -f "$target" ]]; then
                        sed -i.bak "s|${pattern}|${replacement}|g" "$target"
                        rm -f "${target}.bak"
                        echo "[INFO]  codemod: $target ($pattern → $replacement)"
                        applied=$((applied + 1))
                    else
                        echo "[WARN]  codemod target missing: $target" >&2
                    fi
                fi

                i=$((i + 1))
            done
        done < <(manifest_options "$manifest" "$axis")
    done < <(manifest_axes "$manifest")

    echo "[INFO]  Stage F (codemods): applied $applied rules"
}

# stage_npm_deps <manifest_file> <subproject_root>
#
# Removes UNSELECTED options' npm_deps from package.json (dependencies,
# devDependencies, peerDependencies). Deduplicates against SELECTED
# options' npm_deps so shared deps are preserved.
#
# No-op when $root/package.json does not exist (golang, etc.) or when
# no axis declares npm_deps for unselected options.
#
# Requires: jq.
stage_npm_deps() {
    local manifest="$1"
    local root="$2"
    local pkg="$root/package.json"

    if [[ ! -f "$pkg" ]]; then
        echo "[INFO]  Stage G (npm_deps): no package.json at $pkg, skipping"
        return 0
    fi
    command -v jq >/dev/null 2>&1 || {
        echo "[ERROR] Stage G (npm_deps): jq required but not found" >&2
        return 5
    }

    declare -A unselected_deps=()
    declare -A selected_deps=()

    # Pass 1: collect unselected_deps
    while IFS= read -r axis; do
        local selected="${SEL[$axis]:-}"
        while IFS= read -r option; do
            [[ "$option" == "$selected" ]] && continue
            while IFS= read -r dep; do
                [[ -n "$dep" ]] && unselected_deps["$dep"]=1
            done < <(yq -r ".axes.\"$axis\".options.\"$option\".npm_deps // [] | .[]" "$manifest")
        done < <(manifest_options "$manifest" "$axis")
    done < <(manifest_axes "$manifest")

    if (( ${#unselected_deps[@]} == 0 )); then
        echo "[INFO]  Stage G (npm_deps): no npm_deps to remove"
        return 0
    fi

    # Pass 2: collect selected_deps (for dedup)
    while IFS= read -r axis; do
        local selected="${SEL[$axis]:-}"
        [[ -z "$selected" || "$selected" == "none" ]] && continue
        while IFS= read -r dep; do
            [[ -n "$dep" ]] && selected_deps["$dep"]=1
        done < <(yq -r ".axes.\"$axis\".options.\"$selected\".npm_deps // [] | .[]" "$manifest")
    done < <(manifest_axes "$manifest")

    # Build jq filter (single pass)
    local jq_filter='.'
    local removed=0 skipped=0
    for dep in "${!unselected_deps[@]}"; do
        if [[ -n "${selected_deps[$dep]:-}" ]]; then
            skipped=$((skipped + 1))
            continue
        fi
        local esc
        esc=$(printf '%s' "$dep" | jq -R .)
        jq_filter+=" | del(.dependencies[$esc])"
        jq_filter+=" | del(.devDependencies[$esc])"
        jq_filter+=" | del(.peerDependencies[$esc])"
        removed=$((removed + 1))
    done

    if (( removed == 0 )); then
        echo "[INFO]  Stage G (npm_deps): no deletions (all unselected deps shared with selected)"
        return 0
    fi

    local tmp
    tmp=$(mktemp)
    if jq "$jq_filter" "$pkg" > "$tmp" 2>/dev/null; then
        mv "$tmp" "$pkg"
        echo "[INFO]  Stage G (npm_deps): removed $removed dep(s), skipped $skipped (shared with selected)"
        return 0
    else
        rm -f "$tmp"
        echo "[ERROR] Stage G (npm_deps): jq filter failed" >&2
        return 4
    fi
}

# ====================================================================
# SECTION 7: Plan computation and dry-run formatter
# ====================================================================

# compute_plan <manifest_file> <subproject_root>
#
# Populates global arrays:
#   PLAN_DELETIONS  — "axis|option|path" entries (Stage B targets)
#   PLAN_WIRING     — relative paths of files Stage C/D/E will edit
# Pure-read; no filesystem modification.
compute_plan() {
    local manifest="$1"

    PLAN_DELETIONS=()
    PLAN_WIRING=()

    while IFS= read -r axis; do
        local selected="${SEL[$axis]:-}"
        while IFS= read -r option; do
            [[ "$option" == "$selected" ]] && continue
            local path mig sds
            path=$(manifest_option_path "$manifest" "$axis" "$option")
            mig=$(manifest_option_migrations_path "$manifest" "$axis" "$option")
            sds=$(manifest_option_seeds_path "$manifest" "$axis" "$option")
            [[ -n "$path" ]] && PLAN_DELETIONS+=("$axis|$option|$path")
            [[ -n "$mig" ]]  && PLAN_DELETIONS+=("$axis|$option|$mig")
            [[ -n "$sds" ]]  && PLAN_DELETIONS+=("$axis|$option|$sds")
            while IFS= read -r extra; do
                [[ -n "$extra" ]] && PLAN_DELETIONS+=("$axis|$option|$extra")
            done < <(manifest_option_extra_paths "$manifest" "$axis" "$option")
        done < <(manifest_options "$manifest" "$axis")
    done < <(manifest_axes "$manifest")

    while IFS= read -r rel; do
        [[ -n "$rel" ]] && PLAN_WIRING+=("$rel")
    done < <(manifest_wiring_files "$manifest")

    # PLAN_NPM_DEPS: unselected npm_deps that aren't shared with selected.
    PLAN_NPM_DEPS=()
    declare -A _plan_unsel _plan_sel
    while IFS= read -r axis; do
        local selected="${SEL[$axis]:-}"
        while IFS= read -r option; do
            [[ "$option" == "$selected" ]] && continue
            while IFS= read -r dep; do
                [[ -n "$dep" ]] && _plan_unsel["$dep"]=1
            done < <(yq -r ".axes.\"$axis\".options.\"$option\".npm_deps // [] | .[]" "$manifest")
        done < <(manifest_options "$manifest" "$axis")
    done < <(manifest_axes "$manifest")
    while IFS= read -r axis; do
        local selected="${SEL[$axis]:-}"
        [[ -z "$selected" || "$selected" == "none" ]] && continue
        while IFS= read -r dep; do
            [[ -n "$dep" ]] && _plan_sel["$dep"]=1
        done < <(yq -r ".axes.\"$axis\".options.\"$selected\".npm_deps // [] | .[]" "$manifest")
    done < <(manifest_axes "$manifest")
    for dep in "${!_plan_unsel[@]}"; do
        [[ -z "${_plan_sel[$dep]:-}" ]] && PLAN_NPM_DEPS+=("$dep")
    done
    unset _plan_unsel _plan_sel
}

# print_plan <manifest_file> <subproject_root>
#
# Prints a human-readable plan to stdout. Used by --dry-run.
print_plan() {
    local manifest="$1"
    local root="$2"

    echo "[PLAN] Subproject: $SELECTION_SUBPROJECT"
    echo "[PLAN] Manifest:   $manifest"
    echo "[PLAN] Generated:  $(date -u +%Y-%m-%dT%H:%M:%SZ)"
    echo ""
    echo "[PLAN] Selection summary:"
    {
        for axis in "${!SEL[@]}"; do
            printf "  %-15s %s\n" "$axis:" "${SEL[$axis]}"
        done
    } | sort
    echo ""

    echo "=== Stage A — Compose ==="
    echo "  Edit root compose.yml (in-place) via marker filter, subproject=$SELECTION_SUBPROJECT"
    echo ""

    echo "=== Stage B — Adapters ==="
    for entry in "${PLAN_DELETIONS[@]:-}"; do
        [[ -z "$entry" ]] && continue
        IFS='|' read -r axis option path <<< "$entry"
        local full="$root/$path"
        if [[ -e "$full" ]]; then
            printf "  DELETE: %-60s  (axis=%s option=%s)\n" "$full" "$axis" "$option"
        else
            printf "  SKIP   (missing): %-50s  (axis=%s option=%s)\n" "$full" "$axis" "$option"
        fi
    done
    echo "  TOTAL: ${#PLAN_DELETIONS[@]} entries"
    echo ""

    echo "=== Stage C/D/E — Wiring files ==="
    for rel in "${PLAN_WIRING[@]:-}"; do
        [[ -z "$rel" ]] && continue
        echo "  EDIT: $root/$rel"
    done
    echo "  TOTAL: ${#PLAN_WIRING[@]} files"
    echo ""

    echo "=== Stage F — Codemods ==="
    local cm_count=0
    while IFS= read -r axis; do
        local selected="${SEL[$axis]:-}"
        while IFS= read -r option; do
            local n
            n=$(yq -r ".axes.\"$axis\".options.\"$option\".codemods // [] | length" "$manifest")
            local i=0
            while [[ $i -lt $n ]]; do
                local ws wu
                ws=$(yq -r ".axes.\"$axis\".options.\"$option\".codemods[$i].when_selected // false" "$manifest")
                wu=$(yq -r ".axes.\"$axis\".options.\"$option\".codemods[$i].when_unselected // false" "$manifest")
                if [[ "$ws" == "true" && "$selected" == "$option" ]]; then
                    cm_count=$((cm_count + 1))
                elif [[ "$wu" == "true" && "$selected" != "$option" ]]; then
                    cm_count=$((cm_count + 1))
                fi
                i=$((i + 1))
            done
        done < <(manifest_options "$manifest" "$axis")
    done < <(manifest_axes "$manifest")
    if [[ $cm_count -eq 0 ]]; then
        echo "  (no codemods applicable for current selection)"
    else
        echo "  $cm_count codemod rule(s) will fire"
    fi
    echo ""

    echo "=== Stage G — npm_deps ==="
    local npm_count=${#PLAN_NPM_DEPS[@]}
    if (( npm_count == 0 )); then
        echo "  (no npm_deps to remove)"
    else
        for dep in "${PLAN_NPM_DEPS[@]}"; do
            echo "  REMOVE from package.json: $dep"
        done
        echo "  TOTAL: $npm_count dep(s)"
    fi
    echo ""

    echo "=== Stage Verify (will run after apply) ==="
    case "$SELECTION_SUBPROJECT" in
        golang)
            echo "  WILL RUN: cd $root && go build ./..."
            echo "  WILL RUN: cd $root && go mod tidy && git diff --quiet go.mod go.sum"
            ;;
        nextjs)
            echo "  WILL RUN: cd $root && pnpm install --frozen-lockfile=false"
            echo "  WILL RUN: cd $root && pnpm exec tsc --noEmit"
            ;;
        python)
            echo "  WILL RUN: cd $root && uv run python -m compileall -q app"
            echo "  WILL RUN: cd $root && uv run ruff check app tests"
            echo "  WILL RUN: cd $root && uv run mypy app tests/adapters/test_port_parity.py"
            echo "  WILL RUN: cd $root && uv run pytest --collect-only -q"
            ;;
        *)
            echo "  WILL RUN: (subproject-specific verify checks via verify-prune.sh)"
            ;;
    esac
    echo "  WILL RUN: marker pair integrity check across wiring files"
    echo ""
    echo "[PLAN] To apply, re-run without --dry-run."
}

# ====================================================================
# SECTION 8: Main flow
# ====================================================================

main_full_pipeline() {
    # Full mode: --manifest given. Parse → validate → plan →
    # if dry-run print, else apply all stages → verify.
    #
    # .boilerplate.yaml is both the axis declaration AND the selections source
    # (requires a 'selections:' block added by /design-tech-spec).
    local manifest="$MANIFEST_FILE"
    [[ -f "$manifest" ]] || { echo "ERROR: manifest not found: $manifest" >&2; exit 1; }
    parse_selection_file "$manifest" || exit $?
    if [[ ${#SEL[@]} -eq 0 ]]; then
        echo "ERROR: manifest '$manifest' has no 'selections:' block." >&2
        echo "       Add one via /design-tech-spec before pruning." >&2
        exit 1
    fi
    validate_selection "$manifest" || exit $?

    # Resolve folder: read .folder from manifest; fall back to subproject name
    # for manifests that predate the 'folder' key.
    local _folder
    _folder=$(yq -r '.folder // ""' "$manifest")
    SELECTION_FOLDER="${_folder:-$SELECTION_SUBPROJECT}"

    compute_plan "$manifest"

    if [[ $DRY_RUN -eq 1 ]]; then
        print_plan "$manifest" "$SELECTION_FOLDER"
        exit 0
    fi

    echo "[INFO]  Applying prune plan..."
    # Stage A: compose (in-place root compose.yml).
    # Uses SELECTION_SUBPROJECT for boilerplate:subproject= marker matching.
    if [[ -f "compose.yml" ]]; then
        stage_compose "compose.yml" "$SELECTION_SUBPROJECT" "compose.yml" "$manifest" || exit $?
        echo "[INFO]  Stage A (compose): edited compose.yml"
    fi
    # Stage B — filesystem root is SELECTION_FOLDER
    stage_adapters "$manifest" "$SELECTION_FOLDER" || exit $?
    # Stage C/D/E (unified)
    stage_wiring "$manifest" "$SELECTION_FOLDER" || exit $?
    # Stage F
    stage_codemods "$manifest" "$SELECTION_FOLDER" || exit $?
    # Stage G: npm_deps (no-op for subprojects without package.json)
    stage_npm_deps "$manifest" "$SELECTION_FOLDER" || exit $?

    # Stage G.5: settings defaults. After pruning adapters, update *_backend
    # defaults in settings files to match the selected option. Without this,
    # a pruned service may crash at startup because the default references a
    # removed adapter (e.g. storage_backend="local" when only "noop" remains).
    local settings_file="$SELECTION_FOLDER/app/config/settings.py"
    if [[ -f "$settings_file" ]]; then
        local patched=0
        while IFS= read -r axis; do
            local selected="${SEL[$axis]:-}"
            [[ -z "$selected" ]] && continue
            local field="${axis//-/_}_backend"
            if grep -q "${field}:.*str.*=" "$settings_file" 2>/dev/null; then
                local current
                current=$(grep "${field}:" "$settings_file" | sed 's/.*= *"\([^"]*\)".*/\1/')
                if [[ "$current" != "$selected" ]]; then
                    sed -i.bak "s|\(${field}:.*str.*= *\"\)[^\"]*\"|\1${selected}\"|" "$settings_file"
                    rm -f "${settings_file}.bak"
                    echo "[INFO]  settings default: ${field} \"${current}\" → \"${selected}\""
                    patched=$((patched + 1))
                fi
            fi
        done < <(manifest_axes "$manifest")
        echo "[INFO]  Stage G.5 (settings defaults): patched $patched defaults"
    fi

    # Stage H: golang cleanup — remove unused imports then tidy go.mod.
    # goimports must run BEFORE go mod tidy so pruned-adapter imports are gone
    # before tidy tries to validate the module graph.
    if [[ "$SELECTION_SUBPROJECT" == "golang" && -f "$SELECTION_FOLDER/go.mod" ]]; then
        # Locate goimports: prefer $GOPATH/bin, then $HOME/go/bin, then PATH.
        local _goimports=""
        for _candidate in \
            "${GOPATH:-$HOME/go}/bin/goimports" \
            "$HOME/go/bin/goimports" \
            "$(command -v goimports 2>/dev/null || true)"; do
            if [[ -x "$_candidate" ]]; then
                _goimports="$_candidate"
                break
            fi
        done
        if [[ -n "$_goimports" ]]; then
            if "$_goimports" -w "$SELECTION_FOLDER" 2>&1; then
                echo "[INFO]  Stage H (goimports): unused imports removed"
            else
                echo "[WARN]  Stage H (goimports): non-fatal error — continuing" >&2
            fi
        else
            echo "[WARN]  Stage H: goimports not found — skipping unused-import cleanup (install via: go install golang.org/x/tools/cmd/goimports@latest)" >&2
        fi
        if (cd "$SELECTION_FOLDER" && go mod tidy 2>&1); then
            echo "[INFO]  Stage H (go mod tidy): go.mod cleaned"
        else
            echo "[ERROR] Stage H (go mod tidy) failed" >&2
            exit 5
        fi
    fi

    # Stage Verify
    echo "[INFO]  Stage Verify: running scripts/verify-prune.sh --manifest=$manifest"
    local verify_rc=0
    ./scripts/verify-prune.sh --manifest="$manifest" || verify_rc=$?
    if [[ "$verify_rc" -eq 0 ]]; then
        echo "[INFO]  Done. 0 errors."
    elif [[ "$verify_rc" -eq 3 ]]; then
        echo "[ERROR] Constraint violation. See above." >&2
        exit 3
    else
        echo "[ERROR] Verify failed. Inspect via 'git diff' or roll back with 'git reset --hard HEAD'." >&2
        exit 5
    fi
}

main_legacy_compose() {
    # Legacy: --subproject + --selections CSV (+ optional --output)
    if [[ -n "$SELECTIONS" ]]; then
        parse_selection_csv "$SELECTIONS" || exit $?
    fi
    stage_compose "$COMPOSE_FILE" "$SUBPROJECT" "$OUTPUT" "$MANIFEST"
}

# Dispatch
if [[ -n "$MANIFEST_FILE" ]]; then
    main_full_pipeline
else
    main_legacy_compose
fi

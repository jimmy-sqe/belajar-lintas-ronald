#!/usr/bin/env bash
# scripts/verify-prune.sh — Sanity check on a pruned monorepo.
#
# Two modes:
#
# 1) Compose mode (path argument given):
#    Verifies a single pruned compose file:
#      - No remaining `boilerplate:` marker lines
#      - File parses as valid compose (`docker compose config -q`)
#      - No orphan service references
#
# 2) Full mode (--manifest given):
#    Post-prune verification of the entire monorepo:
#      - Marker pair integrity across wiring files (every START has END)
#      - No orphan markers in non-wiring source files
#      - Subproject-specific build check (go build / pnpm tsc / uv compileall+ruff+mypy+collect / gradle compileJava+compileTestJava)
#      - go mod tidy no-op (golang only)
#
# Usage:
#   scripts/verify-prune.sh --manifest=<subproject>   # full mode
#   scripts/verify-prune.sh <pruned-compose-path>     # compose mode
#
# Requirements: bash 4+, yq. Compose mode additionally requires docker.
# Full mode additionally requires go (golang), pnpm (nextjs), uv (python), or
# a JDK + Gradle wrapper (java).

set -euo pipefail

if (( BASH_VERSINFO[0] < 4 )); then
    echo "ERROR: bash 4+ required (found ${BASH_VERSION}); on macOS install via 'brew install bash'" >&2
    exit 1
fi

# Argument parsing — supports --manifest flag + positional compose path
MANIFEST_FILE=""
COMPOSE_PATH=""
while [[ $# -gt 0 ]]; do
    case "$1" in
        --manifest=*) MANIFEST_FILE="${1#*=}" ;;
        --help|-h)
            cat <<EOF
Usage:
  $0 --manifest=<subproject>     # full mode — reads subproject + folder from manifest
  $0 --manifest=<path.yaml>      # full mode — explicit manifest path
  $0 <pruned-compose-path>       # compose mode (single file)
EOF
            exit 0
            ;;
        -*) echo "ERROR: unknown flag: $1" >&2; exit 2 ;;
        *)  COMPOSE_PATH="$1" ;;
    esac
    shift
done

# Normalize --manifest: "golang" → "golang/.boilerplate.yaml"
if [[ -n "$MANIFEST_FILE" && "$MANIFEST_FILE" != */* && "$MANIFEST_FILE" != *.yaml ]]; then
    MANIFEST_FILE="${MANIFEST_FILE}/.boilerplate.yaml"
fi

SUBPROJECT_FILTER=""
FOLDER_PATH=""

if [[ -n "$COMPOSE_PATH" ]]; then
    MODE="compose"
    FILE="$COMPOSE_PATH"
    [[ -f "$FILE" ]] || { echo "ERROR: file not found: $FILE" >&2; exit 3; }
else
    MODE="full"
    if [[ -z "$MANIFEST_FILE" ]]; then
        echo "ERROR: --manifest=<subproject> is required in full mode" >&2
        exit 2
    fi
    [[ -f "$MANIFEST_FILE" ]] || { echo "ERROR: manifest not found: $MANIFEST_FILE" >&2; exit 3; }
    command -v yq >/dev/null 2>&1 || { echo "ERROR: yq is required" >&2; exit 5; }

    SUBPROJECT_FILTER=$(yq -r '.subproject // ""' "$MANIFEST_FILE")
    [[ -n "$SUBPROJECT_FILTER" ]] || {
        echo "ERROR: manifest has no 'subproject' key: $MANIFEST_FILE" >&2
        exit 1
    }
    _folder=$(yq -r '.folder // ""' "$MANIFEST_FILE")
    FOLDER_PATH="${_folder:-$SUBPROJECT_FILTER}"

    # Applies to every stack including Java: persistence=noop maps to the
    # doc-only adapter/persistence/noop package, which neither PersistenceConfig
    # (TodoRepository) nor AuthPersistenceConfig (UserRepository) has an arm for.
    # If sample-app or auth still ship, the context would have no repository bean
    # and fail to start.
    # ── Cross-axis constraint: persistence:noop requires no Repository consumer ──
    _persist=$(yq -r '.selections.persistence // ""' "$MANIFEST_FILE")
    if [[ "$_persist" == "noop" ]]; then
        _sample=$(yq -r '.selections."sample-app" // "none"' "$MANIFEST_FILE")
        _auth=$(yq -r '.selections.auth // "none"' "$MANIFEST_FILE")
        if [[ "$_sample" != "none" || "$_auth" != "none" ]]; then
            echo "FAIL: persistence: noop requires sample-app: none and auth: none (no domain may consume a Repository/user-repo port). Got sample-app=${_sample}, auth=${_auth}." >&2
            exit 3
        fi
    fi

fi

# ---- Full mode: per-subproject dispatch ----

verify_golang_full() {
    local folder="$1"
    local manifest="$folder/.boilerplate.yaml"
    local fail=0

    # Read wiring files from manifest (same source as prune.sh stage_wiring).
    local -a WIRING_FILES=()
    while IFS= read -r f; do
        [[ -z "$f" ]] && continue
        WIRING_FILES+=("$folder/$f")
    done < <(yq -r '.wiring_files // [] | .[]' "$manifest" 2>/dev/null)

    echo "=== 1. Marker pair integrity (wiring files) ==="
    local marker_errors=0
    for f in "${WIRING_FILES[@]}"; do
        [[ -f "$f" ]] || continue
        local starts ends
        starts=$(grep -c 'boilerplate:.*START' "$f" 2>/dev/null || true)
        ends=$(grep -c 'boilerplate:.*END' "$f" 2>/dev/null || true)
        if [[ "$starts" != "$ends" ]]; then
            echo "FAIL: $f has $starts START vs $ends END" >&2
            marker_errors=$((marker_errors + 1))
        fi
    done
    if (( marker_errors == 0 )); then
        echo "OK: all wiring files have balanced START/END"
    else
        fail=1
    fi
    echo

    echo "=== 2. No orphan markers in non-wiring Go files ==="
    local orphan_files=0
    while IFS= read -r f; do
        local skip=0
        for w in "${WIRING_FILES[@]}"; do
            [[ "$f" == "$w" ]] && skip=1 && break
        done
        (( skip == 1 )) && continue
        if grep -q 'boilerplate:' "$f" 2>/dev/null; then
            echo "WARN: orphan marker in non-wiring file: $f" >&2
            orphan_files=$((orphan_files + 1))
        fi
    done < <(find "$folder" -name '*.go' -type f 2>/dev/null)
    if (( orphan_files == 0 )); then
        echo "OK: no orphan markers in non-wiring Go files"
    else
        echo "WARN: $orphan_files file(s) have markers outside the wiring list (review needed)"
    fi
    echo

    echo "=== 3. go build ==="
    if (cd "$folder" && go build ./... 2>&1); then
        echo "OK: go build clean"
    else
        echo "FAIL: go build failed" >&2
        fail=1
    fi
    echo

    echo "=== 4. go mod tidy no-op ==="
    if [[ -f "$folder/go.mod" && -f "$folder/go.sum" ]]; then
        cp "$folder/go.mod" "$folder/go.mod.before"
        cp "$folder/go.sum" "$folder/go.sum.before"
        if (cd "$folder" && go mod tidy 2>&1); then
            if cmp -s "$folder/go.mod" "$folder/go.mod.before" \
                && cmp -s "$folder/go.sum" "$folder/go.sum.before"; then
                echo "OK: go mod tidy produced no diff"
                rm -f "$folder/go.mod.before" "$folder/go.sum.before"
            else
                echo "FAIL: go mod tidy modified go.mod or go.sum (unused deps remaining)" >&2
                diff "$folder/go.mod.before" "$folder/go.mod" >&2 || true
                mv "$folder/go.mod.before" "$folder/go.mod"
                mv "$folder/go.sum.before" "$folder/go.sum"
                fail=1
            fi
        else
            echo "FAIL: go mod tidy errored" >&2
            mv "$folder/go.mod.before" "$folder/go.mod"
            mv "$folder/go.sum.before" "$folder/go.sum"
            fail=1
        fi
    else
        echo "SKIP: go.mod or go.sum not present"
    fi
    echo

    if (( fail )); then
        echo "FAIL: golang full verification did not pass."
        return 5
    fi
    echo "PASS: pruned golang is clean."
    return 0
}

verify_nextjs_full() {
    local folder="$1"
    local manifest="$folder/.boilerplate.yaml"
    local fail=0

    # Read wiring files dynamically from manifest, excluding package.json
    # (handled by stage_npm_deps; not marker-wrapped).
    local -a WIRING_FILES=()
    while IFS= read -r f; do
        [[ -z "$f" ]] && continue
        [[ "$f" == "package.json" ]] && continue
        WIRING_FILES+=("$folder/$f")
    done < <(yq -r '.wiring_files // [] | .[]' "$manifest" 2>/dev/null)

    echo "=== 1. Marker pair integrity (wiring files) ==="
    local marker_errors=0
    for f in "${WIRING_FILES[@]}"; do
        [[ -f "$f" ]] || continue
        local starts ends
        starts=$(grep -c 'boilerplate:.*START' "$f" 2>/dev/null || true)
        ends=$(grep -c 'boilerplate:.*END' "$f" 2>/dev/null || true)
        if [[ "$starts" != "$ends" ]]; then
            echo "FAIL: $f has $starts START vs $ends END" >&2
            marker_errors=$((marker_errors + 1))
        fi
    done
    if (( marker_errors == 0 )); then
        echo "OK: all wiring files have balanced START/END"
    else
        fail=1
    fi
    echo

    echo "=== 2. No orphan markers in non-wiring TS files ==="
    local orphan_count=0
    while IFS= read -r f; do
        local skip=0
        for w in "${WIRING_FILES[@]}"; do
            [[ "$f" == "$w" ]] && skip=1 && break
        done
        (( skip == 1 )) && continue
        if grep -q 'boilerplate:' "$f" 2>/dev/null; then
            echo "WARN: marker in non-wiring file: $f" >&2
            orphan_count=$((orphan_count + 1))
        fi
    done < <(find "$folder/src" -type f \( -name '*.ts' -o -name '*.tsx' \) 2>/dev/null)
    if (( orphan_count == 0 )); then
        echo "OK: no orphan markers"
    else
        echo "WARN: $orphan_count file(s) — review (mocks/handlers.ts is intentional per CLAUDE.md Rule 5)"
    fi
    echo

    echo "=== 3. pnpm install ==="
    if (cd "$folder" && pnpm install --frozen-lockfile=false 2>&1); then
        echo "OK: pnpm install clean"
    else
        echo "FAIL: pnpm install failed" >&2
        fail=1
    fi
    echo

    echo "=== 4. pnpm tsc --noEmit ==="
    if (cd "$folder" && pnpm exec tsc --noEmit 2>&1); then
        echo "OK: TypeScript compile clean"
    else
        echo "FAIL: tsc --noEmit failed" >&2
        fail=1
    fi
    echo

    if (( fail )); then
        echo "FAIL: nextjs full verification did not pass."
        return 5
    fi
    echo "PASS: pruned nextjs is clean."
    return 0
}

verify_qa_test_full() {
    local folder="$1"
    local manifest="$folder/.boilerplate.yaml"
    local fail=0

    # Read wiring files dynamically from manifest.
    local -a WIRING_FILES=()
    while IFS= read -r f; do
        [[ -z "$f" ]] && continue
        WIRING_FILES+=("$folder/$f")
    done < <(yq -r '.wiring_files // [] | .[]' "$manifest" 2>/dev/null)

    echo "=== 1. Marker pair integrity (wiring files) ==="
    local marker_errors=0
    for f in "${WIRING_FILES[@]}"; do
        [[ -f "$f" ]] || continue
        local starts ends
        starts=$(grep -c 'boilerplate:.*START' "$f" 2>/dev/null || true)
        ends=$(grep -c 'boilerplate:.*END' "$f" 2>/dev/null || true)
        if [[ "$starts" != "$ends" ]]; then
            echo "FAIL: $f has $starts START vs $ends END" >&2
            marker_errors=$((marker_errors + 1))
        fi
    done
    if (( marker_errors == 0 )); then
        echo "OK: all wiring files have balanced START/END"
    else
        fail=1
    fi
    echo

    echo "=== 2. No orphan markers in non-wiring Python files ==="
    local orphan_count=0
    while IFS= read -r f; do
        local skip=0
        for w in "${WIRING_FILES[@]}"; do
            [[ "$f" == "$w" ]] && skip=1 && break
        done
        (( skip == 1 )) && continue
        if grep -q 'boilerplate:' "$f" 2>/dev/null; then
            echo "WARN: orphan marker in non-wiring file: $f" >&2
            orphan_count=$((orphan_count + 1))
        fi
    done < <(find "$folder" -type f -name '*.py' -not -path '*/.venv/*' -not -path '*/__pycache__/*' 2>/dev/null)
    if (( orphan_count == 0 )); then
        echo "OK: no orphan markers in non-wiring Python files"
    else
        echo "WARN: $orphan_count file(s) have markers outside the wiring list (review needed)"
    fi
    echo

    echo "=== 3. pip install -e .[dev] ==="
    if (cd "$folder" && pip install -e ".[dev]" -q 2>&1); then
        echo "OK: install clean"
    else
        echo "FAIL: pip install failed" >&2
        fail=1
    fi
    echo

    echo "=== 4. ruff check ==="
    if (cd "$folder" && ruff check . 2>&1); then
        echo "OK: ruff clean"
    else
        echo "FAIL: ruff failed" >&2
        fail=1
    fi
    echo

    echo "=== 5. mypy (strict) ==="
    # Check only the packages that survive the current prune state.
    local mypy_targets=()
    [[ -d "$folder/framework" ]] && mypy_targets+=("framework")
    [[ -d "$folder/pages" ]] && mypy_targets+=("pages")
    [[ -d "$folder/tests" ]] && mypy_targets+=("tests")
    if (cd "$folder" && mypy "${mypy_targets[@]}" 2>&1); then
        echo "OK: mypy clean"
    else
        echo "FAIL: mypy failed (likely a dangling import after prune)" >&2
        fail=1
    fi
    echo

    echo "=== 6. pytest --collect-only ==="
    # Exit code 5 = "no tests collected", which is the EXPECTED outcome when the
    # sample suite is pruned (sample-suite=none). Treat 0 and 5 as success; any
    # other non-zero code means a collection/import error after prune.
    local collect_rc=0
    (cd "$folder" && pytest --collect-only -q 2>&1) || collect_rc=$?
    if [[ "$collect_rc" -eq 0 || "$collect_rc" -eq 5 ]]; then
        echo "OK: collection clean (rc=$collect_rc)"
    else
        echo "FAIL: pytest could not collect (import error after prune, rc=$collect_rc)" >&2
        fail=1
    fi
    echo

    if (( fail )); then
        echo "FAIL: qa-test full verification did not pass."
        return 5
    fi
    echo "PASS: pruned qa-test is clean."
    return 0
}

verify_python_full() {
    local folder="$1"
    local manifest="$folder/.boilerplate.yaml"
    local fail=0

    local -a WIRING_FILES=()
    while IFS= read -r f; do
        [[ -z "$f" ]] && continue
        WIRING_FILES+=("$folder/$f")
    done < <(yq -r '.wiring_files // [] | .[]' "$manifest" 2>/dev/null)

    echo "=== 1. Marker pair integrity (wiring files) ==="
    local marker_errors=0
    for f in "${WIRING_FILES[@]}"; do
        [[ -f "$f" ]] || continue
        local starts ends
        starts=$(grep -c 'boilerplate:.*START' "$f" 2>/dev/null || true)
        ends=$(grep -c 'boilerplate:.*END' "$f" 2>/dev/null || true)
        if [[ "$starts" != "$ends" ]]; then
            echo "FAIL: $f has $starts START vs $ends END" >&2
            marker_errors=$((marker_errors + 1))
        fi
    done
    if (( marker_errors == 0 )); then
        echo "OK: all wiring files have balanced START/END"
    else
        fail=1
    fi
    echo

    echo "=== 2. No orphan markers in non-wiring Python files ==="
    local orphan_count=0
    while IFS= read -r f; do
        local skip=0
        for w in "${WIRING_FILES[@]}"; do
            [[ "$f" == "$w" ]] && skip=1 && break
        done
        (( skip == 1 )) && continue
        if grep -q 'boilerplate:' "$f" 2>/dev/null; then
            echo "WARN: orphan marker in non-wiring file: $f" >&2
            orphan_count=$((orphan_count + 1))
        fi
    done < <(find "$folder" -type f -name '*.py' -not -path '*/.venv/*' -not -path '*/__pycache__/*' 2>/dev/null)
    if (( orphan_count == 0 )); then
        echo "OK: no orphan markers in non-wiring Python files"
    else
        echo "WARN: $orphan_count file(s) have markers outside the wiring list (review needed)"
    fi
    echo

    echo "=== 3. compileall ==="
    if (cd "$folder" && uv run python -m compileall -q app 2>&1); then
        echo "OK: compileall clean"
    else
        echo "FAIL: compileall failed (syntax error after prune)" >&2
        fail=1
    fi
    echo

    echo "=== 4. ruff check ==="
    if (cd "$folder" && uv run ruff check app tests 2>&1); then
        echo "OK: ruff clean"
    else
        echo "FAIL: ruff check failed" >&2
        fail=1
    fi
    echo

    echo "=== 5. mypy ==="
    local _mypy_targets="app"
    [[ -f "$folder/tests/adapters/test_port_parity.py" ]] && _mypy_targets="$_mypy_targets tests/adapters/test_port_parity.py"
    if (cd "$folder" && uv run mypy $_mypy_targets 2>&1); then
        echo "OK: mypy clean"
    else
        echo "FAIL: mypy failed (likely a dangling import after prune)" >&2
        fail=1
    fi
    echo

    echo "=== 6. pytest --collect-only ==="
    local collect_rc=0
    (cd "$folder" && uv run pytest --collect-only -q 2>&1) || collect_rc=$?
    if [[ "$collect_rc" -eq 0 || "$collect_rc" -eq 5 ]]; then
        echo "OK: collection clean (rc=$collect_rc)"
    else
        echo "FAIL: pytest could not collect (import error after prune, rc=$collect_rc)" >&2
        fail=1
    fi
    echo

    if (( fail )); then
        echo "FAIL: python full verification did not pass."
        return 5
    fi
    echo "PASS: pruned python is clean."
    return 0
}

verify_java_full() {
    local folder="$1"
    local manifest="$folder/.boilerplate.yaml"
    local fail=0

    local -a WIRING_FILES=()
    while IFS= read -r f; do
        [[ -z "$f" ]] && continue
        WIRING_FILES+=("$folder/$f")
    done < <(yq -r '.wiring_files // [] | .[]' "$manifest" 2>/dev/null)

    echo "=== 1. Marker pair integrity (wiring files) ==="
    local marker_errors=0
    for f in "${WIRING_FILES[@]}"; do
        [[ -f "$f" ]] || continue
        local starts ends
        starts=$(grep -c 'boilerplate:.*START' "$f" 2>/dev/null || true)
        ends=$(grep -c 'boilerplate:.*END' "$f" 2>/dev/null || true)
        if [[ "$starts" != "$ends" ]]; then
            echo "FAIL: $f has $starts START vs $ends END" >&2
            marker_errors=$((marker_errors + 1))
        fi
    done
    if (( marker_errors == 0 )); then
        echo "OK: all wiring files have balanced START/END"
    else
        fail=1
    fi
    echo

    echo "=== 2. No orphan markers in non-wiring Java files ==="
    local orphan_count=0
    while IFS= read -r f; do
        local skip=0
        for w in "${WIRING_FILES[@]}"; do
            [[ "$f" == "$w" ]] && skip=1 && break
        done
        (( skip == 1 )) && continue
        if grep -q 'boilerplate:' "$f" 2>/dev/null; then
            echo "WARN: orphan marker in non-wiring file: $f" >&2
            orphan_count=$((orphan_count + 1))
        fi
    done < <(find "$folder" -type f -name '*.java' -not -path '*/build/*' 2>/dev/null)
    if (( orphan_count == 0 )); then
        echo "OK: no orphan markers in non-wiring Java files"
    else
        echo "WARN: $orphan_count file(s) have markers outside the wiring list (review needed)"
    fi
    echo

    echo "=== 3. Marker axes are declared in the manifest ==="
    # Every axis referenced by a boilerplate marker MUST be a key under .axes,
    # or the pruner (which is manifest-driven) can neither delete its files nor
    # strip its marker blocks — the option becomes permanently baked in. This is
    # the check that catches "wired everywhere but never registered" gaps.
    local -A MANIFEST_AXES=()
    while IFS= read -r ax; do
        [[ -z "$ax" ]] && continue
        MANIFEST_AXES["$ax"]=1
    done < <(yq -r '.axes // {} | keys | .[]' "$manifest" 2>/dev/null)
    local axis_errors=0
    while IFS= read -r ax; do
        [[ -z "$ax" ]] && continue
        if [[ -z "${MANIFEST_AXES[$ax]:-}" ]]; then
            echo "FAIL: marker references axis '$ax' but it is not declared in .axes (prune cannot strip it)" >&2
            axis_errors=$((axis_errors + 1))
        fi
    done < <(grep -rho --include='*.java' --include='*.gradle' --include='*.yml' --include='*.yaml' \
                 --exclude='.boilerplate.yaml' --exclude-dir=build 'axis=[a-zA-Z0-9_-]*' "$folder" 2>/dev/null \
             | sed 's/axis=//' | sort -u)
    if (( axis_errors == 0 )); then
        echo "OK: every marker axis is declared in the manifest"
    else
        fail=1
    fi
    echo

    echo "=== 4. compileJava ==="
    # Compiles main sources (triggers generateProto when rpc=grpc survives).
    # A dangling import left by a bad prune fails here.
    if (cd "$folder" && ./gradlew --quiet compileJava 2>&1); then
        echo "OK: compileJava clean"
    else
        echo "FAIL: compileJava failed (likely a dangling import after prune)" >&2
        fail=1
    fi
    echo

    echo "=== 5. compileTestJava ==="
    # Compiles test sources so pruned/orphaned tests referencing deleted
    # adapters or generated stubs surface as failures.
    if (cd "$folder" && ./gradlew --quiet compileTestJava 2>&1); then
        echo "OK: compileTestJava clean"
    else
        echo "FAIL: compileTestJava failed (test references a pruned symbol)" >&2
        fail=1
    fi
    echo

    if (( fail )); then
        echo "FAIL: java full verification did not pass."
        return 5
    fi
    echo "PASS: pruned java is clean."
    return 0
}

if [[ "$MODE" == "full" ]]; then
    # Resolve stack from folder contents when subproject name is non-canonical
    # (e.g. --add-backend=python:ai-engine creates subproject=ai-engine, stack=python).
    _resolve_stack() {
        case "$1" in
            golang|nextjs|python|qa-test|java) echo "$1"; return ;;
        esac
        local folder="$2"
        if [[ -f "$folder/go.mod" ]]; then echo "golang"
        elif [[ -f "$folder/build.gradle" || -f "$folder/build.gradle.kts" || -f "$folder/settings.gradle" ]]; then echo "java"
        elif [[ -f "$folder/pyproject.toml" || -f "$folder/setup.py" ]]; then echo "python"
        elif [[ -f "$folder/package.json" ]] && grep -q '"next"' "$folder/package.json" 2>/dev/null; then echo "nextjs"
        elif [[ -f "$folder/requirements.txt" ]] && grep -q "playwright" "$folder/requirements.txt" 2>/dev/null; then echo "qa-test"
        else echo ""; fi
    }
    RESOLVED_STACK=$(_resolve_stack "$SUBPROJECT_FILTER" "$FOLDER_PATH")
    if [[ -z "$RESOLVED_STACK" ]]; then
        echo "ERROR: cannot determine stack for subproject '$SUBPROJECT_FILTER' at '$FOLDER_PATH' (no go.mod, build.gradle, pyproject.toml, or package.json found)" >&2
        exit 1
    fi
    case "$RESOLVED_STACK" in
        golang)  verify_golang_full "$FOLDER_PATH" ;;
        nextjs)  verify_nextjs_full "$FOLDER_PATH" ;;
        python)  verify_python_full "$FOLDER_PATH" ;;
        qa-test) verify_qa_test_full "$FOLDER_PATH" ;;
        java)    verify_java_full "$FOLDER_PATH" ;;
    esac
    exit $?
fi

# ---- Compose mode (legacy single-file verification) ----

fail=0

echo "=== 1. No remaining marker lines ==="
markers=$(grep -c 'boilerplate:' "$FILE" || true)
if [[ "$markers" -gt 0 ]]; then
    echo "FAIL: $markers marker line(s) remain in $FILE" >&2
    grep -n 'boilerplate:' "$FILE" >&2
    fail=1
else
    echo "OK: no marker lines remain"
fi

echo
echo "=== 2. Compose syntax valid ==="
if docker compose -f "$FILE" config -q 2>/tmp/compose-config.err; then
    echo "OK: docker compose config -q passes"
else
    echo "FAIL: docker compose config -q rejects the file"
    cat /tmp/compose-config.err >&2
    rm -f /tmp/compose-config.err
    fail=1
fi

echo
echo "=== 3. depends_on integrity ==="
command -v yq >/dev/null 2>&1 || { echo "SKIP: yq not available — depends_on integrity check skipped"; echo; if (( fail )); then echo "FAIL: see issues above."; exit 1; fi; echo "PASS: pruned compose is clean."; exit 0; }

services=$(yq -r '.services | keys | .[]' "$FILE" 2>/dev/null || true)
declare -A SVC_SET=()
for s in $services; do SVC_SET["$s"]=1; done

orphan=0
# Handle both depends_on forms: map ({db: {...}}) yields keys; list (- db) yields scalars directly.
while IFS= read -r dep_target; do
    [[ -z "$dep_target" ]] && continue
    if [[ -z "${SVC_SET[$dep_target]:-}" ]]; then
        echo "FAIL: depends_on references non-existent service: $dep_target" >&2
        orphan=1
    fi
done < <(yq -r '
  .services[] | .depends_on // null |
  (select(. != null) |
    (select((. | type) == "!!map") | keys | .[]),
    (select((. | type) == "!!seq") | .[]))
' "$FILE" 2>/dev/null || true)

if (( orphan == 0 )); then
    echo "OK: no orphan depends_on references"
else
    fail=1
fi

echo
if (( fail )); then
    echo "FAIL: verification did not pass."
    exit 1
fi
echo "PASS: pruned compose is clean."
exit 0

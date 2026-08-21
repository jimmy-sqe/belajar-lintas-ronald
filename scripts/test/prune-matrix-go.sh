#!/usr/bin/env bash
# Prune-matrix for the go stack. Each row: clean archive → write selections →
# prune → assert exit code. Needs bash>=4 + yq + go on PATH.
set -uo pipefail
REPO="$(cd "$(dirname "$0")/../.." && pwd)"

run_row() {
    local name="$1" want="$2" persist="$3" cache="$4" auth="$5" sample="$6"
    local work; work="$(mktemp -d)"
    git -C "$REPO" archive HEAD | tar -x -C "$work"
    yq -i ".subproject = \"golang\" | .folder = \"golang\"" "$work/golang/.boilerplate.yaml"
    yq -i ".selections = {\"persistence\": \"$persist\", \"cache\": \"$cache\", \"auth\": \"$auth\", \"sample-app\": \"$sample\", \"storage\": \"noop\", \"messaging\": \"noop\", \"notification\": \"noop\", \"secrets\": \"env\", \"observability\": \"noop\", \"api-doc\": \"none\", \"jobs\": \"noop\", \"container\": \"dockerfile\", \"rpc\": \"none\"}" "$work/golang/.boilerplate.yaml"
    ( cd "$work" && bash scripts/prune.sh --manifest=golang/.boilerplate.yaml ) >/dev/null 2>&1
    local got=$?
    rm -rf "$work"
    if [[ "$got" == "$want" ]]; then echo "PASS: $name (exit $got)"; else echo "FAIL: $name — expected $want, got $got"; return 1; fi
}

fail=0
run_row "default-postgres-todo-jwt" 0 postgres redis    jwt  todo-app || fail=1
run_row "mysql-todo-authnone"       0 mysql    inmemory none todo-app || fail=1
run_row "mongodb-todo-authnone"     0 mongodb  inmemory none todo-app || fail=1
run_row "stateless-noop"            0 noop     noop     none none     || fail=1
run_row "constraint-noop-with-todo" 3 noop     inmemory none todo-app || fail=1
exit $fail

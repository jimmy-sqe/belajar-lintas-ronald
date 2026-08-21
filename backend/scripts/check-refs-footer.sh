#!/usr/bin/env bash
# check-refs-footer.sh — commit-msg hook helper
#
# Enforces "Refs: <TICKET-ID>" footer for feat/fix commits.
# Per docs/CONTRIBUTING.md: feat and fix MUST include the footer.
#
# Invoked by pre-commit framework with the path to .git/COMMIT_EDITMSG.

set -euo pipefail

msg_file="${1:-.git/COMMIT_EDITMSG}"
[[ -f "$msg_file" ]] || { echo "check-refs-footer: $msg_file not found" >&2; exit 1; }

msg=$(cat "$msg_file")

# Strip leading comments (git sometimes inserts them)
subject=$(echo "$msg" | grep -v '^#' | head -1)

# Extract type from "<type>(scope): subject" or "<type>: subject"
type=$(echo "$subject" | sed -nE 's/^([a-z]+)(\(.+\))?!?:.*/\1/p')

case "$type" in
    feat|fix)
        if ! echo "$msg" | grep -qE '^Refs: [A-Z]+-[0-9]+'; then
            echo "ERROR: $type commit requires 'Refs: <TICKET-ID>' footer (e.g., 'Refs: SDLC-123')" >&2
            echo "       See docs/CONTRIBUTING.md for the full convention." >&2
            exit 1
        fi
        ;;
    "")
        # Not a Conventional Commit — conventional-pre-commit handles that error.
        :
        ;;
esac

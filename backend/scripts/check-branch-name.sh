#!/usr/bin/env bash
# check-branch-name.sh — pre-push hook helper
#
# Enforces branch naming: <type>/<TICKET-ID>-<slug>
# See docs/CONTRIBUTING.md.

set -euo pipefail

branch=$(git rev-parse --abbrev-ref HEAD)

# Allow main and HEAD (detached) without checks.
case "$branch" in
    main|HEAD) exit 0 ;;
esac

# Format: <type>/<TICKET-ID>-<slug>, max 50 chars total.
if [[ ! "$branch" =~ ^(feature|fix|chore|docs|experiment)/[A-Z]+-[0-9]+-[a-z0-9-]+$ ]]; then
    echo "ERROR: Branch '$branch' does not match required format." >&2
    echo "       Expected: <type>/<TICKET-ID>-<slug>" >&2
    echo "       Types:    feature | fix | chore | docs | experiment" >&2
    echo "       Example:  feature/SDLC-123-payment-flow" >&2
    exit 1
fi

if [[ ${#branch} -gt 50 ]]; then
    echo "ERROR: Branch name exceeds 50 characters." >&2
    exit 1
fi

exit 0

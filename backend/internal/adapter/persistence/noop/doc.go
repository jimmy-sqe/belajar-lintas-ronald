// Package noop is the no-op persistence option: the service has no database, so no
// repository is wired. Selecting persistence=noop requires sample-app=none and
// auth=none (enforced by scripts/verify-prune.sh).
package noop

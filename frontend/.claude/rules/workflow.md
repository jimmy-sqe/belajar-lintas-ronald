# Workflow Rules — Commit & PR

Detail for the commit/PR directives summarized in `CLAUDE.md`.

## Commit & PR

- Conventional Commits: `<type>(<scope>): <subject>` with body
  explaining **why**, not **what**.
- `feat` and `fix` REQUIRE `Refs: <TICKET-ID>` footer.
- Branch: `<type>/<TICKET-ID>-<slug>` (≤ 50 chars, kebab-case slug).
- NEVER `--no-verify`. Pre-commit Husky hooks run lint + format check.

→ Detail: `docs/CONTRIBUTING.md`

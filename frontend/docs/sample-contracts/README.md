# Sample Auth Contracts (Reference)

These YAML files document the **expected backend contract shape** for each auth-model adapter under `src/services/auth/<option>/`.

**They are documentation only.** Kubb does not consume them. The source of truth for FE adapter types is `src/services/auth/<option>/types.ts` (hand-written).

Use these YAMLs to:
- Communicate the required backend API to your backend team
- Author MSW handlers in `src/mocks/auth/<option>/handlers.ts`
- Onboard new contributors to the auth model

If your real backend contract diverges from a sample (e.g. different field names or path prefixes), update **both** the active adapter's `types.ts` AND the corresponding YAML reference to match.

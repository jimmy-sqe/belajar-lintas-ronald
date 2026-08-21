/**
 * Post-login default landing route + redirect target when an
 * unauthenticated user hits a protected page.
 *
 * Plugin-codemod hint:
 *   - sample-app axis includes `todo-app` → '/todos' (boilerplate default)
 *   - sample-app axis excludes `todo-app`  → '/' (or another sample's
 *     landing path picked by the plugin)
 *
 * The value is intentionally NOT marker-wrapped so the export survives
 * regardless of axis selection; the plugin replaces only the literal
 * string at scaffold time. Consumers (e.g. useLoginPage) always have a
 * value to import.
 */
export const defaultBaseUrlPage = '/todos';

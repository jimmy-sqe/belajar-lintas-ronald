/**
 * Server barrel. Ships pointing at the noop register. The pruner
 * repoints the literal below to the selected option. Kept separate from
 * index.ts so server startup never loads a browser-only SDK.
 */
export { register } from './noop/register';

/**
 * Client barrel. Ships pointing at the noop provider. The pruner
 * repoints the literal below to the selected option (stage_codemods,
 * when_selected: true) and deletes the other option folders. App code imports
 * `ObservabilityProvider` from '@/services/observability' only.
 */
export { default as ObservabilityProvider } from './noop/provider';
export * from './types';

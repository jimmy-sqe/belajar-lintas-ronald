/**
 * Next.js instrumentation entrypoint. Delegates to the selected
 * observability adapter's server-side register().
 */
export { register } from '@/services/observability/server';

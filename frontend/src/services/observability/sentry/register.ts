import * as Sentry from '@sentry/nextjs';

/** Sentry server-side init, invoked from instrumentation.ts. */
export async function register(): Promise<void> {
  const dsn = process.env.NEXT_PUBLIC_SENTRY_DSN;
  if (!dsn) return;
  Sentry.init({ dsn, tracesSampleRate: 1.0 });
}

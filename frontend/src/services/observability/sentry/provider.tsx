'use client';

import { useEffect } from 'react';
import * as Sentry from '@sentry/nextjs';

import type { ObservabilityProviderProps } from '../types';

/** Sentry client-side init (errors + tracing). */
export default function ObservabilityProvider({ children }: ObservabilityProviderProps) {
  useEffect(() => {
    const dsn = process.env.NEXT_PUBLIC_SENTRY_DSN;
    if (!dsn) return;
    Sentry.init({ dsn, tracesSampleRate: 1.0 });
  }, []);

  return <>{children}</>;
}

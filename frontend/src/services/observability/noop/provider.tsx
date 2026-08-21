'use client';

import type { ObservabilityProviderProps } from '../types';

/** No-op provider: renders children with no instrumentation. */
export default function ObservabilityProvider({ children }: ObservabilityProviderProps) {
  return <>{children}</>;
}

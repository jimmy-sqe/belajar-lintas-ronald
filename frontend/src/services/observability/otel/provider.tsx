'use client';

import type { ObservabilityProviderProps } from '../types';

/**
 * OTel here is server-centric (traces via @vercel/otel in register.ts).
 * The client surface is an explicit passthrough.
 */
export default function ObservabilityProvider({ children }: ObservabilityProviderProps) {
  return <>{children}</>;
}

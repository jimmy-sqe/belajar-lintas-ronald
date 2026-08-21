'use client';

import { useEffect } from 'react';
import { AwsRum, type AwsRumConfig } from 'aws-rum-web';

import type { ObservabilityProviderProps } from '../types';

/** AWS CloudWatch RUM client provider. Init must never break the app. */
export default function ObservabilityProvider({ children }: ObservabilityProviderProps) {
  useEffect(() => {
    const appId = process.env.NEXT_PUBLIC_RUM_APP_ID;
    const region = process.env.NEXT_PUBLIC_RUM_REGION;
    if (!appId || !region) return;

    try {
      const config: AwsRumConfig = {
        sessionSampleRate: 1,
        identityPoolId: process.env.NEXT_PUBLIC_RUM_IDENTITY_POOL_ID,
        endpoint: process.env.NEXT_PUBLIC_RUM_ENDPOINT,
        telemetries: ['errors', 'http', 'performance'],
        allowCookies: true,
        enableXRay: false
      };
      new AwsRum(appId, '1.0.0', region, config);
    } catch {
      // Swallow RUM init failures — observability must not crash the app.
    }
  }, []);

  return <>{children}</>;
}

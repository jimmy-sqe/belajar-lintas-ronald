export type RenderingMode = 'ssr' | 'csr';

export const RENDERING_MODE: RenderingMode =
  (process.env.NEXT_PUBLIC_RENDERING_MODE as RenderingMode) || 'csr';

export const isSSR = () => RENDERING_MODE === 'ssr';

// boilerplate:axis=auth option=jwt-refresh START
/**
 * Feature flag: enable/disable automatic refresh token logic.
 * Set NEXT_PUBLIC_ENABLE_REFRESH_TOKEN=true in .env to enable.
 */
export const isRefreshTokenEnabled = () => process.env.NEXT_PUBLIC_ENABLE_REFRESH_TOKEN === 'true';
// boilerplate:axis=auth option=jwt-refresh END

export const isDev = process.env.NODE_ENV === 'development';

export const appName = process.env.NEXT_PUBLIC_APP_NAME || 'SQEBoilerplate';

export const apiBaseUrl = (
  process.env.NEXT_PUBLIC_BACKEND_URL || 'https://api.example.com'
).replace(/\/$/, '');

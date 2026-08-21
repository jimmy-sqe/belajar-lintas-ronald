import { isDev, isSSR } from '@/config/environment';
import type { FetchResponse, FetchResult, RequestConfig } from './parser';
// boilerplate:axis=auth option=jwt-refresh START
import { withRefreshQueue } from '@/services/auth/jwt-refresh/fetcher-middleware';
// boilerplate:axis=auth option=jwt-refresh END

export type { Client, RequestConfig, ResponseErrorConfig } from './parser';

/**
 * Client factory that returns the appropriate client based on the rendering mode.
 *
 * - SSR mode: uses `fetcher-ssr.ts` (proxies through server action)
 * - CSR mode: uses `fetcher-csr.ts` (direct Axios with interceptors)
 *
 * Kubb config points `client.importPath` to this file so generated hooks
 * work transparently in both modes.
 */
const baseFetch = async <TData, _TError = unknown, TVariables = unknown>(
  config: RequestConfig<TVariables>
): Promise<FetchResponse<TData>> => {
  let result: FetchResult<TData>;
  if (isSSR()) {
    const { serverFetch } = await import('./fetcher-ssr');
    result = await serverFetch(config);
  } else {
    const { clientFetch } = await import('./fetcher-csr');
    result = await clientFetch(config);
  }

  if (!result.success) {
    if (isDev && result.message !== 'canceled') {
      console.error(isSSR() ? 'SSR Fetch error:' : 'CSR Fetch error:', result, config);
    }

    // Reconstruct the error locally so `cause` is accessible to the caller.
    throw new Error(result.message, { cause: result.cause });
  }

  return result?.data;
};

// boilerplate:axis=auth option=jwt-refresh START
export const customFetch = withRefreshQueue(baseFetch);
// boilerplate:axis=auth option=jwt-refresh END

export default customFetch;

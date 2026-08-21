'use server';

import axios from 'axios';
import auth from '@/services/auth';
import type { RequestConfig } from '@/services/fetcher/parser';
import {
  parseResponse,
  refreshTokenFailedResponse,
  type FetchResult
} from '@/services/fetcher/parser';
import {
  createRefreshQueueState,
  processQueue
} from '@/services/auth/jwt-refresh/fetcher-middleware';
import { isRefreshTokenEnabled } from '@/config/environment';
import { makeRequest } from './axios';

/**
 * Global state for token refresh management
 * - isRefreshing: indicates if a token refresh is in progress
 * - requestQueue: stores pending requests during token refresh
 */
const tokenRefreshState = createRefreshQueueState();

/**
 * Server Action: Server Fetch
 *
 * Proxies an API request through the server action:
 * 1. Reads the encrypted session cookie
 * 2. Injects the access token into the Authorization header
 * 3. Forwards the request to the backend
 * 4. On 401, refreshes the token and retries once
 * 5. Returns a serializable response to the client
 */
export async function serverFetch<TData = unknown>(
  config: RequestConfig
): Promise<FetchResult<TData>> {
  const session = await auth.getSession();

  // Check if refresh token is expired — session is fully dead, no point sending the request
  if (
    !!session &&
    isRefreshTokenEnabled() &&
    session.expiresAt !== null &&
    Date.now() > session.expiresAt
  ) {
    await auth.signOut();
    return { success: false, message: 'Session expired', cause: null };
  }

  // Strip client-only non-serializable fields (e.g. AbortSignal passed by React Query).
  // Accessing `signal.aborted` on a client reference from the server throws a RSC error.

  const { signal: _signal, cancelToken: _cancelToken, ...safeConfig } = config;

  try {
    const response = await makeRequest<TData>(safeConfig, session?.token || '');
    return { success: true, data: parseResponse(response) };
  } catch (error) {
    // On 401, refresh the token and retry once
    if (
      axios.isAxiosError(error) &&
      (error.response?.status === 401 || error.response?.data?.code === 1005) &&
      !!session
    ) {
      if (isRefreshTokenEnabled()) {
        // Queue request if refresh is already in progress
        if (tokenRefreshState.isRefreshing) {
          return new Promise<FetchResult<TData>>((resolve, reject) => {
            tokenRefreshState.requestQueue.push({
              resolve: resolve as (value: FetchResult<unknown>) => void,
              reject,
              config: safeConfig
            });
          });
        }

        tokenRefreshState.isRefreshing = true;
        let newSession;
        try {
          newSession = await auth.refresh?.(session);
        } catch (refreshError) {
          await processQueue(tokenRefreshState, refreshError);
          return refreshTokenFailedResponse<TData>();
        }

        if (!newSession) {
          await processQueue(tokenRefreshState, new Error('Session expired and refresh failed.'));
          return refreshTokenFailedResponse<TData>();
        }

        await processQueue(tokenRefreshState, null, newSession.token);
        const retryResponse = await makeRequest<TData>(safeConfig, newSession.token || '');
        return { success: true, data: parseResponse(retryResponse) };
      } else {
        await auth.signOut();
        return refreshTokenFailedResponse<TData>();
      }
    }

    if (axios.isAxiosError(error) && error.response) {
      return {
        success: false,
        message: error.message,
        cause: parseResponse(error.response, { statusText: error.message })
      };
    }

    return {
      success: false,
      message: error instanceof Error ? error.message : String(error),
      cause: error instanceof Error ? (error.cause ?? null) : null
    };
  }
}

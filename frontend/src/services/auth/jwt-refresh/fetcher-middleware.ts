/**
 * Higher-order fetcher wrapper that queues concurrent requests during
 * a refresh-in-flight and replays them with the new access token.
 * Only mounted when auth=jwt-refresh.
 */

import type { NextRequest } from 'next/server';
import { makeRequest } from '@/services/fetcher/axios';
import { parseResponse, type FetchResult, type RequestConfig } from '@/services/fetcher/parser';
import { isAccessExpired } from './session';
import { getStoredSession, setStoredSession } from './storage';
import { refresh } from './auth';

export type QueueItem = {
  resolve: (value: FetchResult<unknown>) => void;
  reject: (reason?: unknown) => void;
  config: RequestConfig;
};

export type RefreshQueueState = {
  isRefreshing: boolean;
  requestQueue: QueueItem[];
};

export const createRefreshQueueState = (): RefreshQueueState => ({
  isRefreshing: false,
  requestQueue: []
});

/**
 * Process all queued requests after token refresh.
 * Resolves each queued request with the new token, or rejects on error.
 * Resets the queue state when done.
 *
 * @param state - The shared token refresh state object
 * @param error - Optional error if token refresh failed
 * @param accessToken - New access token to use for retried requests
 */
export const processQueue = async (
  state: RefreshQueueState,
  error: unknown = null,
  accessToken?: string
): Promise<void> => {
  for (const request of state.requestQueue) {
    if (error) {
      request.reject(error);
      continue;
    }
    try {
      const response = await makeRequest(request.config, accessToken || '');
      request.resolve({ success: true, data: parseResponse(response) });
    } catch (err) {
      request.reject(err);
    }
  }

  // Reset queue state
  state.requestQueue = [];
  state.isRefreshing = false;
};

// Module-level shared refresh state (singleton per adapter instance)
const _state = createRefreshQueueState();

/**
 * Higher-order function: wraps a fetcher with refresh-queue behaviour.
 * Concurrent 401/expiry situations are queued and replayed with the new token.
 *
 * @param fetcherFn - The base fetcher to wrap
 * @returns A fetcher with transparent token refresh and queuing
 */
export function withRefreshQueue<T extends (...args: any[]) => Promise<any>>(fetcherFn: T): T {
  return (async (...args: Parameters<T>) => {
    const config: RequestConfig = args[0] ?? {};

    // If already refreshing, queue this request
    if (_state.isRefreshing) {
      return new Promise<FetchResult<unknown>>((resolve, reject) => {
        _state.requestQueue.push({ resolve, reject, config });
      });
    }

    try {
      return await fetcherFn(...args);
    } catch (err: any) {
      const status = err?.response?.status ?? err?.status;
      if (status !== 401) throw err;

      // First 401 — start refresh
      if (_state.isRefreshing) {
        return new Promise<FetchResult<unknown>>((resolve, reject) => {
          _state.requestQueue.push({ resolve, reject, config });
        });
      }

      _state.isRefreshing = true;
      try {
        const stored = await getStoredSession();
        // normalizeSession is not needed here; refresh() handles persistence
        const session = stored
          ? {
              token: stored.access_token,
              user: stored.user,
              permissions: stored.permissions,
              expiresAt: null
            }
          : null;

        if (!session) throw new Error('No session for refresh');
        const newSession = await refresh(session as any);
        await processQueue(_state, null, newSession.token);
        // Retry original request with new token
        return await fetcherFn(...args);
      } catch (refreshErr) {
        await processQueue(_state, refreshErr);
        throw refreshErr;
      }
    }
  }) as T;
}

/**
 * Best-effort SSR middleware guard: reads the stored session and refreshes
 * the access token if it has expired. Errors are swallowed — the guard is
 * advisory only; the downstream API will return 401 if truly invalid.
 *
 * Called from `src/middleware.ts` before the request is forwarded.
 *
 * @param _req - The incoming Next.js request (unused; cookie access via next/headers)
 */
export async function tryRefreshIfExpired(_req: NextRequest): Promise<void> {
  try {
    const stored = await getStoredSession();
    if (!stored) return;
    if (!isAccessExpired(stored)) return;

    const session = {
      token: stored.access_token,
      user: stored.user,
      permissions: stored.permissions,
      expiresAt: null
    };
    await refresh(session as any);
  } catch {
    // Best-effort — swallow all errors
  }
}

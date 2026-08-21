import type { AxiosRequestConfig, AxiosResponse } from 'axios';
import axios from 'axios';
import { apiBaseUrl } from '@/config/environment';
import { buildFetchConfig } from './parser';

/**
 * Plain Axios client used by server actions (proxy, auth) and other API requests.
 *
 * - No auth interceptors (tokens are injected manually by callers)
 * - No localStorage access (runs on the server)
 * - Uses NEXT_PUBLIC_BACKEND_URL as the base URL
 */
const api = axios.create({
  baseURL: apiBaseUrl,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
});

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const makeRequest = async <TData, TError = unknown, TVariables = unknown>(
  config: AxiosRequestConfig<TVariables>,
  token?: string
): Promise<AxiosResponse<TData>> => {
  return api.request<TData, AxiosResponse<TData>, TVariables>(
    buildFetchConfig({
      ...config,
      headers: {
        ...config.headers,
        Authorization: token ? `Bearer ${token}` : undefined
      }
    })
  );
};

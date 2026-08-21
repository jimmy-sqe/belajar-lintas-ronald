import type {
  AxiosError,
  AxiosRequestConfig,
  AxiosResponse,
  InternalAxiosRequestConfig
} from 'axios';

export type FetchResponse<TData = unknown> = {
  data: TData;
  status: number;
  statusText: string;
  headers: Record<string, string>;
};

export type FetchResult<TData = unknown> =
  | { success: true; data: FetchResponse<TData> }
  | { success: false; message: string; cause: unknown };

export type RequestConfig<TVariables = unknown> = AxiosRequestConfig<TVariables>;
export type ResponseErrorConfig<TError = unknown> = AxiosError<TError>;
type RequestConfigWithError<TVariables = unknown, TError = unknown> = RequestConfig<TVariables> & {
  _error?: TError;
};
export type Client = <TData, TError = unknown, TVariables = unknown>(
  config: RequestConfigWithError<TVariables, TError>
) => Promise<AxiosResponse<TData>>;

export type QueueItem = {
  resolve: (value: InternalAxiosRequestConfig | PromiseLike<InternalAxiosRequestConfig>) => void;
  reject: (reason?: Error | unknown) => void;
  config: InternalAxiosRequestConfig;
};

export function parseResponse<TData>(
  response: AxiosResponse<TData>,
  overrides?: Partial<FetchResponse<TData>>
): FetchResponse<TData> {
  const headers: Record<string, string> = {};
  if (response.headers) {
    for (const [key, value] of Object.entries(response.headers)) {
      if (typeof value === 'string') {
        headers[key] = value;
      }
    }
  }

  return {
    data: response.data,
    status: response.status,
    statusText: response.statusText,
    headers,
    ...overrides
  };
}

export function buildFetchConfig<TVariables>(
  config: RequestConfig<TVariables>
): RequestConfig<TVariables> {
  return {
    ...config,
    headers: {
      ...config.headers,
      'x-utc-offset': (-new Date().getTimezoneOffset()).toString()
    }
  };
}

export function refreshTokenFailedResponse<TData>(): FetchResult<TData> {
  return {
    success: false,
    message: 'Session expired and refresh failed.',
    cause: {
      status: 401,
      statusText: 'Session expired and refresh failed',
      data: { msg: 'Session expired and refresh failed' }
    }
  };
}

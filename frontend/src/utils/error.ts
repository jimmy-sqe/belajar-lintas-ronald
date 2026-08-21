import type { AxiosError } from 'axios';

interface IErrorResponseData {
  code: string;
  msg: string;
}

export function errorResponseHandler(err: unknown) {
  const error = err as AxiosError;
  const errorResponse = error.response?.data as IErrorResponseData;
  const messageArr = errorResponse?.msg?.split('|');
  return {
    message: messageArr || ['Something went wrong'],
    code: errorResponse?.code || 'UNKNOWN_ERROR'
  };
}

export const extractErrorApi = <T>(error: AxiosError) => {
  const data = error?.response?.data as T;
  return data;
};

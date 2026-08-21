import { useCallback } from 'react';
import { useToasterProvider } from '@squantumengine/horizon';
import type { Query } from '@tanstack/react-query';
import type { AxiosError } from 'axios';
import type { HttpErrorResponse } from '@/services/types';
import { extractErrorApi } from '@/utils/error';

export type ToastErrorHandler = <
  TError = unknown,
  TQueryFnData = unknown,
  TData = unknown,
  TQueryKey extends readonly unknown[] = readonly unknown[]
>(
  error: TError,
  query?: Query<TQueryFnData, TError, TData, TQueryKey>
) => void;

const useToastError = (): { onError: ToastErrorHandler } => {
  const { open } = useToasterProvider();

  const handleCopy = useCallback(
    async (id: string) => {
      await navigator.clipboard.writeText(id);
      open({
        id: 'copy-toaster',
        label: 'ID berhasil disalin'
      });
    },
    [open]
  );

  const onError: ToastErrorHandler = useCallback(
    (error, query) => {
      if (query?.meta?.skipToast) return;
      const axiosError = error as AxiosError;
      const statusCode = axiosError?.response?.status ?? 500;
      const data = extractErrorApi<HttpErrorResponse>(axiosError);
      const time = Date.now().toString();
      open({
        id: `error-toaster-${time}`,
        label: (
          <span>
            {statusCode >= 500 ? 'Sistem sedang terkendala. Silakan coba lagi nanti.' : data?.message}
          </span>
        ),
        variant: 'secondary',
        buttonLabel: <span className="font-bold">{statusCode >= 500 ? 'Salin ID' : 'OK'}</span>,
        onActionClick: () => {
          if (statusCode >= 500) {
            handleCopy(`error-toaster-${time}`);
          }
        }
      });
    },
    [handleCopy, open]
  );

  return {
    onError
  };
};

export default useToastError;

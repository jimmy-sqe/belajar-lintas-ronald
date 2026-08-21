'use client';

import { Suspense, use } from 'react';
import { handlers as mockHandlers } from '@/mocks/handlers';
import { handlersPromise } from '@/utils/msw/handlers';

const mockingEnabledPromise =
  typeof window !== 'undefined' && process.env.NEXT_PUBLIC_API_MOCKING === 'enabled'
    ? import('@/utils/msw/browser').then(async ({ worker }) => {
        await worker.start({
          onUnhandledRequest(request, print) {
            if (request.url.includes('_next')) {
              return;
            }
            print.warning();
          }
        });

        try {
          const openApiHandlers = await handlersPromise;
          worker.use(...openApiHandlers, ...mockHandlers);
        } catch (error) {
          console.error('Failed to load OpenAPI handlers in MSWProvider:', error);
        }
      })
    : Promise.resolve();

export function MSWProvider({
  children
}: Readonly<{
  children: React.ReactNode;
}>) {
  // If MSW is enabled, we need to wait for the worker to start,
  // so we wrap the children in a Suspense boundary until it's ready.
  return (
    <Suspense fallback={null}>
      <MSWProviderWrapper>{children}</MSWProviderWrapper>
    </Suspense>
  );
}

function MSWProviderWrapper({
  children
}: Readonly<{
  children: React.ReactNode;
}>) {
  use(mockingEnabledPromise);
  return children;
}

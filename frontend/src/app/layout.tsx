import { Suspense } from 'react';
import { ToasterProvider } from '@squantumengine/horizon';
import type { Metadata } from 'next';
import localFont from 'next/font/local';
import { z } from 'zod';

import './globals.css';
import '@squantumengine/horizon/lib/index.css';

import LoadingDialog from '@/components/loadingDialog';
import ReactQueryProvider from '@/components/providers/reactQuery';
import SessionProvider from '@/components/providers/session';
// boilerplate:axis=mocking option=msw START
import { MSWProvider } from '@/utils/msw/msw-provider';
// boilerplate:axis=mocking option=msw END
import { ObservabilityProvider } from '@/services/observability';
import { customZodErrorMap } from '@/utils/zod';

z.setErrorMap(customZodErrorMap);

const geistSans = localFont({
  src: './fonts/GeistVF.woff',
  variable: '--font-geist-sans',
  weight: '100 900'
});
const geistMono = localFont({
  src: './fonts/GeistMonoVF.woff',
  variable: '--font-geist-mono',
  weight: '100 900'
});

export const metadata: Metadata = {
  title: 'frontend-belajar-lintas-ronald',
  description: 'Greenfield Next.js boilerplate (App Router + React 19 + Horizon UI + Kubb)'
};

export default function RootLayout({
  children
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} theme-hz-blue antialiased`}
        suppressHydrationWarning>
        <Suspense>
          <ObservabilityProvider>
            <MSWProvider>
              <ToasterProvider
                settings={{
                  placement: 'top',
                  timeout: 3
                }}>
                <ReactQueryProvider>
                  <SessionProvider>{children}</SessionProvider>
                </ReactQueryProvider>
              </ToasterProvider>
              <LoadingDialog />
            </MSWProvider>
          </ObservabilityProvider>
        </Suspense>
      </body>
    </html>
  );
}

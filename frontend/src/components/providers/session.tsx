import { dehydrate, HydrationBoundary } from '@tanstack/react-query';
import { getQueryClient } from '@/services/query-client';
import auth from '@/services/auth';
import { SESSION_QUERY_KEY } from '@/hooks/auth/useSession';

export default async function SessionProvider({ children }: { children: React.ReactNode }) {
  const queryClient = getQueryClient();
  const session = await auth.getSession();

  if (session) {
    queryClient.setQueryData(SESSION_QUERY_KEY, session);
  }

  return <HydrationBoundary state={dehydrate(queryClient)}>{children}</HydrationBoundary>;
}

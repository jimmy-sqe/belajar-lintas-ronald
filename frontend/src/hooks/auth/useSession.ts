'use client';

import { useQuery } from '@tanstack/react-query';
import auth from '@/services/auth';

export const SESSION_QUERY_KEY = ['session'] as const;

const useSession = () => {
  return useQuery({
    queryKey: SESSION_QUERY_KEY,
    queryFn: () => auth.getSession(),
    staleTime: Infinity,
    refetchOnWindowFocus: false
  });
};

export default useSession;

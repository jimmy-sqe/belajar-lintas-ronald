'use client';

import { useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import auth from '@/services/auth';
import useToastError from '../ui/useToastError';
import { SESSION_QUERY_KEY } from './useSession';

const useLogout = () => {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { onError } = useToastError();

  const handleLogout = async () => {
    try {
      await auth.signOut();
      queryClient.setQueryData(SESSION_QUERY_KEY, null);
      router.replace('/login');
    } catch (error) {
      onError(error as Error);
    }
  };

  return {
    logout: handleLogout
  };
};

export default useLogout;

'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import auth from '@/services/auth';
import { isSSR } from '@/config/environment';
import { defaultBaseUrlPage } from '@/static/route';

const AuthenticationRouteCSR = ({ children }: { children: React.ReactNode }) => {
  const router = useRouter();
  const [isChecking, setIsChecking] = useState(true);

  useEffect(() => {
    if (isSSR()) {
      setIsChecking(false);
      return;
    }
    const checkAuth = async () => {
      const session = await auth.getSession();
      if (session) {
        return router.replace(defaultBaseUrlPage);
      }
      setIsChecking(false);
    };

    checkAuth();
  }, [router]);

  if (isChecking) return null;

  return children;
};

export default AuthenticationRouteCSR;

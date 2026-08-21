'use client';

import { useEffect, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import auth from '@/services/auth';
import { isSSR } from '@/config/environment';

const ProtectAuthenticationRouteCSR = ({ children }: { children: React.ReactNode }) => {
  const router = useRouter();
  const firstRender = useRef(false);
  const [isChecking, setIsChecking] = useState(true);

  useEffect(() => {
    if (isSSR()) {
      setIsChecking(false);
      return;
    }
    if (firstRender.current) return;
    firstRender.current = true;

    const checkAuth = async () => {
      const session = await auth.getSession();
      if (!session) {
        const pathname = location.pathname;
        return router.replace('/login?callbackUrl=' + encodeURIComponent(pathname));
      }
      setIsChecking(false);
    };

    checkAuth();
  }, [router]);

  if (isChecking) return null;

  return children;
};

export default ProtectAuthenticationRouteCSR;

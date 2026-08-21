import { headers } from 'next/headers';
import { redirect } from 'next/navigation';
import auth from '@/services/auth';
import { isSSR } from '@/config/environment';
import ProtectAuthenticationRouteCSR from './protectAuthenticationRouteCSR';

const ProtectAuthenticationRoute = async ({ children }: { children: React.ReactNode }) => {
  if (isSSR()) {
    const session = await auth.getSession();
    if (!session) {
      const headersList = await headers();
      const pathname = headersList.get('x-pathname') ?? '/';
      redirect('/login?callbackUrl=' + encodeURIComponent(pathname));
    }

    return children;
  }

  return <ProtectAuthenticationRouteCSR>{children}</ProtectAuthenticationRouteCSR>;
};

export default ProtectAuthenticationRoute;

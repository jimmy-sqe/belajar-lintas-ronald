import { redirect } from 'next/navigation';
import auth from '@/services/auth';
import { isSSR } from '@/config/environment';
import { defaultBaseUrlPage } from '@/static/route';
import AuthenticationRouteCSR from './authenticationRouteCSR';

const AuthenticationRoute = async ({ children }: { children: React.ReactNode }) => {
  if (isSSR()) {
    const session = await auth.getSession();
    if (session) {
      redirect(defaultBaseUrlPage);
    }

    return children;
  }

  return <AuthenticationRouteCSR>{children}</AuthenticationRouteCSR>;
};

export default AuthenticationRoute;

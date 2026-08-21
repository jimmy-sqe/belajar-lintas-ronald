import type React from 'react';
import AuthenticationRoute from '@/components/auth/authenticationRoute';

const AuthLayout = ({ children }: { children: React.ReactNode }) => {
  return <AuthenticationRoute>{children}</AuthenticationRoute>;
};

export default AuthLayout;

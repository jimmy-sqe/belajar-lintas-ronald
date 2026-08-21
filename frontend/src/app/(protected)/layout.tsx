import ProtectAuthentication from '@/components/auth/protectAuthenticationRoute';

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  return <ProtectAuthentication>{children}</ProtectAuthentication>;
}

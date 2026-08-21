'use client';

import { Sidebar } from './Sidebar';
import { Topbar } from './Topbar';

interface AppShellProps {
  children: React.ReactNode;
  /** Override breadcrumbs (default: ["Workspace", "Todos"]). */
  crumbs?: string[];
}

export function AppShell({ children, crumbs }: AppShellProps) {
  return (
    <div className="screen" style={{ position: 'relative' }}>
      <div className="app">
        <Sidebar />
        <div className="main">
          <Topbar crumbs={crumbs} />
          <div className="main__body">
            <div className="main__inner">{children}</div>
          </div>
        </div>
      </div>
    </div>
  );
}

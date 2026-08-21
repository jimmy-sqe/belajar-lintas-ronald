'use client';

import { useRouter, usePathname } from 'next/navigation';
import { Icon, useToaster } from '@squantumengine/horizon';
import useSession from '@/hooks/auth/useSession';
import auth from '@/services/auth';
import { useTodos } from '@/services/todos';

function getInitials(name?: string): string {
  if (!name) return '?';
  const parts = name.trim().split(/\s+/).filter(Boolean);
  if (parts.length === 0) return '?';
  if (parts.length === 1) return parts[0]!.slice(0, 2).toUpperCase();
  return (parts[0]![0]! + parts[parts.length - 1]![0]!).toUpperCase();
}

export function Sidebar() {
  const router = useRouter();
  const pathname = usePathname();
  const { data: session } = useSession();
  const toaster = useToaster();
  // Sidebar count uses a separate query key (page_size: 1) so it does
  // not interfere with the list page's useTodos cache.
  const { data: countQuery } = useTodos({ page: 1, page_size: 1 });
  const total = countQuery?.pagination?.total_data;

  const onSignOut = async () => {
    try {
      await auth.signOut();
      router.push('/login');
    } catch {
      // Spec §10: surface failure to user so they can retry; do NOT
      // navigate — keep them on the protected route.
      // The horizon OpenToaster union forces buttonLabel/onActionClick to be
      // `never` when omitted; cast keeps the call site ergonomic.
      toaster.open({
        id: `signout-${Date.now()}`,
        label: 'Sign-out failed, please retry',
        variant: 'secondary',
      } as Parameters<typeof toaster.open>[0]);
    }
  };

  const isTodosActive = pathname?.startsWith('/todos');

  return (
    <aside className="sb">
      <div className="sb__brand">
        <div className="sb__mark">
          <Icon name="check" />
        </div>
        <div>
          <div className="sb__brand-name">TodoApp</div>
          <div className="sb__brand-sub">internal · v0.1</div>
        </div>
      </div>

      <div className="sb__section-label">Workspace</div>
      <div className="sb__nav">
        <button
          className={'sb__item' + (isTodosActive ? ' is-active' : '')}
          onClick={() => router.push('/todos')}>
          <Icon name="list-check" />
          Todos
          {typeof total === 'number' && <span className="sb__count">{total}</span>}
        </button>
      </div>

      <div className="sb__section-label">Help</div>
      <div className="sb__nav">
        <button className="sb__item" disabled title="Not yet available">
          <Icon name="book" />
          README
        </button>
        <button className="sb__item" disabled title="Not yet available">
          <Icon name="keyboard" />
          Shortcuts
        </button>
      </div>

      <div className="sb__spacer" />

      <div className="sb__user">
        <div className="sb__avatar">{getInitials(session?.user?.name)}</div>
        <div style={{ minWidth: 0, flex: 1 }}>
          <div className="sb__user-name">{session?.user?.name ?? 'Anonymous'}</div>
          <div className="sb__user-email">{session?.user?.email ?? ''}</div>
        </div>
        <button
          type="button"
          onClick={onSignOut}
          aria-label="Sign out"
          style={{ background: 'transparent', border: 0, cursor: 'pointer' }}>
          <Icon name="sign-out" />
        </button>
      </div>
    </aside>
  );
}

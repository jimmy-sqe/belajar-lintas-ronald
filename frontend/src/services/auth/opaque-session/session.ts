import type { Session } from '../types';
import type { OpaqueStoredSession } from './types';
import { getStoredSession } from './storage';

export function normalizeSession(stored: OpaqueStoredSession): Session {
  return {
    token: stored.session_token,
    user: { id: stored.user.id, name: stored.user.name, email: stored.user.email },
    permissions: stored.permissions,
    expiresAt: null // server-side TTL; client cannot decode
  };
}

export async function getSession(): Promise<Session | null> {
  const stored = await getStoredSession();
  return stored ? normalizeSession(stored) : null;
}

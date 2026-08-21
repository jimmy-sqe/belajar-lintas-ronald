import customFetch from '@/services/fetcher';
import type { Session, SignInRequest } from '../types';
import type { JwtSignInResponse, JwtRefreshResponse, JwtStoredSession } from './types';
import { setStoredSession, clearStoredSession, getStoredSession } from './storage';
import { decodeExpiry } from './session';

export async function signIn(request: SignInRequest): Promise<Session> {
  const res = await customFetch<JwtSignInResponse>({
    url: '/v1/auth/login',
    method: 'POST',
    data: request
  });
  const stored: JwtStoredSession = {
    ...res.data.data,
    permissions: res.data.data.permissions ?? []
  };
  await setStoredSession(stored);
  return normalizeSession(stored);
}

export async function signOut(): Promise<void> {
  await customFetch({ url: '/v1/auth/logout', method: 'POST' }).catch(() => undefined);
  await clearStoredSession();
}

export async function refresh(_session: Session): Promise<Session> {
  const stored = await getStoredSession();
  if (!stored) throw new Error('refresh called without a stored session');
  const res = await customFetch<JwtRefreshResponse>({
    url: '/v1/auth/renew',
    method: 'POST',
    data: { refresh_token: stored.refresh_token }
  });
  const updated: JwtStoredSession = {
    ...stored,
    access_token: res.data.data.access_token,
    access_token_expires_in_sec: res.data.data.access_token_expires_in_sec,
    issued_at: res.data.data.issued_at
  };
  await setStoredSession(updated);
  return normalizeSession(updated);
}

export function normalizeSession(stored: JwtStoredSession): Session {
  return {
    token: stored.access_token,
    user: { id: stored.user.id, name: stored.user.name, email: stored.user.email },
    permissions: stored.permissions,
    expiresAt: decodeExpiry(stored)
  };
}

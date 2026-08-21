import customFetch from '@/services/fetcher';
import type { Session, SignInRequest, ChangePasswordRequest } from '../types';
import type { OpaqueSignInResponse, OpaqueStoredSession } from './types';
import { setStoredSession, clearStoredSession } from './storage';
import { normalizeSession } from './session';

export async function signIn(request: SignInRequest): Promise<Session> {
  const res = await customFetch<OpaqueSignInResponse>({
    url: '/auth/sign-in',
    method: 'POST',
    data: request
  });
  const stored: OpaqueStoredSession = {
    session_token: res.data.data.session_token,
    user: res.data.data.user,
    permissions: res.data.data.permissions,
    requires_password_change: res.data.data.requires_password_change
  };
  await setStoredSession(stored);
  return normalizeSession(stored);
}

export async function signOut(): Promise<void> {
  await customFetch({ url: '/auth/sign-out', method: 'POST' }).catch(() => undefined);
  await clearStoredSession();
}

export async function changePassword(req: ChangePasswordRequest): Promise<void> {
  await customFetch({ url: '/auth/change-password', method: 'POST', data: req });
}

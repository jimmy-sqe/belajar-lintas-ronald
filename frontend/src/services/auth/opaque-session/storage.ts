import { isSSR } from '@/config/environment';
import type { OpaqueStoredSession } from './types';

const SESSION_KEY = 'sqe.session';

async function csrSet(stored: OpaqueStoredSession): Promise<void> {
  localStorage.setItem(SESSION_KEY, JSON.stringify(stored));
}

async function csrGet(): Promise<OpaqueStoredSession | null> {
  if (typeof localStorage === 'undefined') return null;
  const raw = localStorage.getItem(SESSION_KEY);
  return raw ? (JSON.parse(raw) as OpaqueStoredSession) : null;
}

async function csrClear(): Promise<void> {
  if (typeof localStorage === 'undefined') return;
  localStorage.removeItem(SESSION_KEY);
}

async function ssrSet(stored: OpaqueStoredSession): Promise<void> {
  const { setSessionCookie } = await import('./storage-ssr');
  await setSessionCookie(stored);
}

async function ssrGet(): Promise<OpaqueStoredSession | null> {
  const { getSessionCookie } = await import('./storage-ssr');
  return getSessionCookie();
}

async function ssrClear(): Promise<void> {
  const { clearSessionCookie } = await import('./storage-ssr');
  await clearSessionCookie();
}

export async function setStoredSession(stored: OpaqueStoredSession): Promise<void> {
  return isSSR() ? ssrSet(stored) : csrSet(stored);
}

export async function getStoredSession(): Promise<OpaqueStoredSession | null> {
  return isSSR() ? ssrGet() : csrGet();
}

export async function clearStoredSession(): Promise<void> {
  return isSSR() ? ssrClear() : csrClear();
}

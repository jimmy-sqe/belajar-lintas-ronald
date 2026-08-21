import { isSSR } from '@/config/environment';
import type { JwtStoredSession } from './types';

const SESSION_KEY = 'sqe.session';

// ============================================================
// CSR — encrypted localStorage
// ============================================================

async function csrSet(stored: JwtStoredSession): Promise<void> {
  const { setEncryptedItem } = await import('@/utils/crypto/storage');
  await setEncryptedItem(SESSION_KEY, JSON.stringify(stored));
}

async function csrGet(): Promise<JwtStoredSession | null> {
  const { getDecryptedItem } = await import('@/utils/crypto/storage');
  const raw = await getDecryptedItem(SESSION_KEY);
  return raw ? (JSON.parse(raw) as JwtStoredSession) : null;
}

async function csrClear(): Promise<void> {
  const { removeEncryptedItem } = await import('@/utils/crypto/storage');
  removeEncryptedItem(SESSION_KEY);
}

// ============================================================
// SSR — encrypted httpOnly cookie via server action
// ============================================================

async function ssrSet(stored: JwtStoredSession): Promise<void> {
  const { setSessionCookie } = await import('./storage-ssr');
  await setSessionCookie(stored);
}

async function ssrGet(): Promise<JwtStoredSession | null> {
  const { getSessionCookie } = await import('./storage-ssr');
  return getSessionCookie();
}

async function ssrClear(): Promise<void> {
  const { clearSessionCookie } = await import('./storage-ssr');
  await clearSessionCookie();
}

// ============================================================
// Public API — branches by rendering mode
// ============================================================

export async function setStoredSession(stored: JwtStoredSession): Promise<void> {
  return isSSR() ? ssrSet(stored) : csrSet(stored);
}

export async function getStoredSession(): Promise<JwtStoredSession | null> {
  return isSSR() ? ssrGet() : csrGet();
}

export async function clearStoredSession(): Promise<void> {
  return isSSR() ? ssrClear() : csrClear();
}

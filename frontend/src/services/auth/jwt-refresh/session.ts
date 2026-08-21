import type { Session } from '../types';
import type { JwtStoredSession } from './types';
import { getStoredSession } from './storage';
import { normalizeSession } from './auth';

/** Buffer (seconds) before actual expiry to treat token as expired. */
export const EXPIRY_BUFFER_SEC = 30;

/**
 * Compute access-token expiry in ms-since-epoch.
 * Returns null if the stored session lacks expiry fields.
 */
export function decodeExpiry(stored: JwtStoredSession): number | null {
  if (!stored.access_token || !stored.access_token_expires_in_sec) return null;
  return stored.issued_at + (stored.access_token_expires_in_sec - EXPIRY_BUFFER_SEC) * 1000;
}

export function isAccessExpired(stored: JwtStoredSession): boolean {
  const expiresAt = decodeExpiry(stored);
  if (expiresAt === null) return true;
  return Date.now() > expiresAt;
}

export function isRefreshExpired(stored: JwtStoredSession): boolean {
  if (!stored.refresh_token || !stored.refresh_token_expires_in_sec) return true;
  const expiresAt =
    stored.issued_at + (stored.refresh_token_expires_in_sec - EXPIRY_BUFFER_SEC) * 1000;
  return Date.now() > expiresAt;
}

/** Read the session from disk and normalize for app consumption. */
export async function getSession(): Promise<Session | null> {
  const stored = await getStoredSession();
  if (!stored) return null;
  return normalizeSession(stored);
}

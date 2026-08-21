'use server';

import { cookies, headers } from 'next/headers';
import { sessionKeyEnum } from '@/static/session';
import type { JwtStoredSession } from './types';
import { decryptSession, encryptSession } from '@/utils/crypto/session';

const COOKIE_NAME = sessionKeyEnum.SESSION;

/**
 * Encrypt session data and set it as an httpOnly cookie
 */
export async function setSessionCookie(data: JwtStoredSession): Promise<void> {
  const maxAge = data.refresh_token_expires_in_sec ?? 0;
  const jwt = await encryptSession(data as Parameters<typeof encryptSession>[0], maxAge);

  const cookieStore = await cookies();
  cookieStore.set(COOKIE_NAME, jwt, {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict',
    maxAge,
    path: '/'
  });
}

/**
 * Read and decrypt the session cookie.
 *
 * Middleware may inject two special request headers before SSR runs:
 *   x-session-cleared: 1       → treat session as absent (cookie was invalidated)
 *   x-session-jwt: <jwt>       → use this refreshed JWT instead of the stale cookie
 *
 * This lets the page see the correct session state in the same request in
 * which middleware refreshed or cleared the cookie.
 */
export async function getSessionCookie(): Promise<JwtStoredSession | null> {
  try {
    const headersList = await headers();

    // Middleware signalled the cookie was cleared this request
    if (headersList.get('x-session-cleared') === '1') return null;

    // Middleware refreshed the token this request — use the new JWT directly
    const refreshedJwt = headersList.get('x-session-jwt');
    if (refreshedJwt) {
      return (await decryptSession(refreshedJwt)) as JwtStoredSession | null;
    }

    // Normal path: read from the cookie
    const cookieStore = await cookies();
    const cookie = cookieStore.get(COOKIE_NAME);
    if (!cookie?.value) return null;

    return (await decryptSession(cookie.value)) as JwtStoredSession | null;
  } catch {
    return null;
  }
}

/**
 * Clear the session cookie
 */
export async function clearSessionCookie(): Promise<void> {
  const cookieStore = await cookies();
  cookieStore.delete(COOKIE_NAME);
}

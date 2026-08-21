'use server';

import { cookies } from 'next/headers';
import type { OpaqueStoredSession } from './types';

const COOKIE_NAME = 'sqe.session';
const MAX_AGE_SEC = 24 * 60 * 60; // 24h — match COSMOS session TTL

export async function setSessionCookie(stored: OpaqueStoredSession): Promise<void> {
  const cookieStore = await cookies();
  cookieStore.set(COOKIE_NAME, JSON.stringify(stored), {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'lax',
    maxAge: MAX_AGE_SEC,
    path: '/'
  });
}

export async function getSessionCookie(): Promise<OpaqueStoredSession | null> {
  const cookieStore = await cookies();
  const raw = cookieStore.get(COOKIE_NAME)?.value;
  if (!raw) return null;
  try {
    return JSON.parse(raw) as OpaqueStoredSession;
  } catch {
    return null;
  }
}

export async function clearSessionCookie(): Promise<void> {
  const cookieStore = await cookies();
  cookieStore.delete(COOKIE_NAME);
}

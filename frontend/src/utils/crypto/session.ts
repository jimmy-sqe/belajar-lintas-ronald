/**
 * Edge-compatible session crypto utilities.
 * No 'use server', no next/headers — safe to import in middleware (Edge Runtime).
 */

import { EncryptJWT, jwtDecrypt } from 'jose';
import type { JwtStoredSession } from '@/services/auth/jwt-refresh/types';

export const EXPIRY_BUFFER_SEC = 60;

function getSecret(): Uint8Array {
  const secret = process.env.AUTH_COOKIE_SECRET;
  if (!secret || secret.length < 32)
    throw new Error('AUTH_COOKIE_SECRET must be at least 32 characters');
  return new TextEncoder().encode(secret);
}

export async function encryptSession(data: JwtStoredSession, maxAge: number): Promise<string> {
  return new EncryptJWT(data as unknown as Record<string, unknown>)
    .setProtectedHeader({ alg: 'dir', enc: 'A256GCM' })
    .setIssuedAt()
    .setExpirationTime(`${maxAge}s`)
    .encrypt(getSecret());
}

export async function decryptSession(jwt: string): Promise<JwtStoredSession | null> {
  try {
    const { payload } = await jwtDecrypt(jwt, getSecret());
    return payload as unknown as JwtStoredSession;
  } catch {
    return null;
  }
}

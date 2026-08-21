/**
 * Encrypted localStorage wrapper using Web Crypto API (AES-GCM).
 *
 * Provides the same get/set/remove interface as localStorage but encrypts
 * values before storing and decrypts them on read.
 */

const ALGORITHM = 'AES-GCM';
const IV_LENGTH = 12;

/**
 * Get or derive the encryption key from the environment variable.
 * Falls back to a default key in development for convenience.
 */
async function getCryptoKey(): Promise<CryptoKey> {
  const keyMaterial =
    process.env.NEXT_PUBLIC_ENCRYPTION_KEY || 'default-dev-key-do-not-use-in-prod!';

  const encoder = new TextEncoder();
  const keyData = encoder.encode(keyMaterial);

  // Hash the key material to ensure exactly 256 bits
  const hash = await crypto.subtle.digest('SHA-256', keyData);

  return crypto.subtle.importKey('raw', hash, { name: ALGORITHM }, false, ['encrypt', 'decrypt']);
}

/**
 * Encrypt a string value and store it in localStorage
 */
export async function setEncryptedItem(key: string, value: string): Promise<void> {
  if (typeof window === 'undefined') return;

  try {
    const cryptoKey = await getCryptoKey();
    const iv = crypto.getRandomValues(new Uint8Array(IV_LENGTH));
    const encoder = new TextEncoder();
    const data = encoder.encode(value);

    const encrypted = await crypto.subtle.encrypt({ name: ALGORITHM, iv }, cryptoKey, data);

    // Store as base64: iv + ciphertext
    const combined = new Uint8Array(iv.length + encrypted.byteLength);
    combined.set(iv, 0);
    combined.set(new Uint8Array(encrypted), iv.length);

    // Chunk to avoid stack overflow and poor perf when spreading large Uint8Arrays
    let binary = '';
    const chunkSize = 0x8000;
    for (let i = 0; i < combined.length; i += chunkSize) {
      binary += String.fromCharCode(...combined.subarray(i, i + chunkSize));
    }
    localStorage.setItem(key, btoa(binary));
  } catch (error) {
    console.error('Failed to encrypt item:', error);
    // Fallback to plain storage in case of crypto failure
    localStorage.setItem(key, value);
  }
}

/**
 * Read and decrypt a value from localStorage
 */
export async function getDecryptedItem(key: string): Promise<string | null> {
  if (typeof window === 'undefined') return null;

  const stored = localStorage.getItem(key);
  if (!stored) return null;

  try {
    const cryptoKey = await getCryptoKey();
    const combined = Uint8Array.from(atob(stored), (c) => c.charCodeAt(0));

    const iv = combined.slice(0, IV_LENGTH);
    const ciphertext = combined.slice(IV_LENGTH);

    const decrypted = await crypto.subtle.decrypt({ name: ALGORITHM, iv }, cryptoKey, ciphertext);

    return new TextDecoder().decode(decrypted);
  } catch {
    // If decryption fails, it might be a plain-text value (migration case)
    // Return it as-is so the user can still use it until they re-login
    return stored;
  }
}

/**
 * Remove an item from localStorage
 */
export function removeEncryptedItem(key: string): void {
  if (typeof window === 'undefined') return;
  localStorage.removeItem(key);
}

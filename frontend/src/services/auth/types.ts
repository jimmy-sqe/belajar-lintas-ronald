/**
 * Shared types for the auth subsystem. All auth-model adapters
 * (jwt-refresh, opaque-session, oauth2-oidc, none) implement this
 * contract so app code can consume one stable interface regardless
 * of which model is active.
 */

export interface Session {
  /** Bearer token to inject into Authorization header. */
  token: string;
  user: SessionUser;
  /** RBAC strings (empty array if model has no concept). */
  permissions: string[];
  /**
   * ms-since-epoch. Only set if the model can decode expiry
   * client-side (JWT). `null` means "trust 401 from server".
   */
  expiresAt: number | null;
}

export interface SessionUser {
  id: string;
  name?: string;
  email?: string;
}

export type AuthErrorCode =
  | 'INVALID_CREDENTIALS'
  | 'SESSION_EXPIRED'
  | 'FORBIDDEN'
  | 'NETWORK'
  | 'UNKNOWN';

export class AuthError extends Error {
  constructor(
    public readonly code: AuthErrorCode,
    message: string,
    public readonly cause?: unknown
  ) {
    super(message);
    this.name = 'AuthError';
  }
}

export interface SignInRequest {
  email: string;
  password: string;
}

export interface ChangePasswordRequest {
  currentPassword: string;
  newPassword: string;
}

export interface ResetPasswordRequest {
  email: string;
}

/**
 * Adapter contract. Each `src/services/auth/<option>/index.ts`
 * default-exports an object satisfying this interface.
 */
export interface AuthAdapter {
  signIn(request: SignInRequest): Promise<Session>;
  signOut(): Promise<void>;
  getSession(): Promise<Session | null>;
  /** Only for models that support refresh (jwt-refresh). */
  refresh?: (session: Session) => Promise<Session>;
  changePassword?: (request: ChangePasswordRequest) => Promise<void>;
  /** Only for models that support self-service reset (jwt-refresh). */
  resetPassword?: (request: ResetPasswordRequest) => Promise<void>;
}

/**
 * Opaque session paradigm. Backend returns an opaque session_token in a
 * standard envelope. Client cannot decode expiry; must rely on 401 from
 * server. No refresh — re-login on expiry.
 *
 * Reference contract: docs/sample-contracts/opaque-session.openapi.yaml
 */

export interface OpaqueSignInResponse {
  success: true;
  code: number;
  data: {
    session_token: string;
    user: { id: string; name?: string; email?: string };
    permissions: string[];
    requires_password_change?: boolean;
  };
  timestamp?: string;
}

export interface OpaqueChangePasswordResponse {
  success: true;
  code: number;
  data: { changed_at: string };
  timestamp?: string;
}

/** Persistent session shape on disk. */
export interface OpaqueStoredSession {
  session_token: string;
  user: { id: string; name?: string; email?: string };
  permissions: string[];
  requires_password_change?: boolean;
}

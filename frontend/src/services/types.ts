/**
 * Generic backend error envelope. Hand-written; not derived from
 * kubb output. Used across auth and non-auth API responses.
 */
export interface HttpErrorResponse {
  success: false;
  code: number;
  message: string;
  timestamp: string;
  metadata?: Record<string, unknown>;
}

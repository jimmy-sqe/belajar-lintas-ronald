import type { AuthAdapter } from '../types';

function notImplemented(): never {
  throw new Error(
    'OAuth2/OIDC auth model is not implemented in v1.0.0. ' +
      'Implement using NextAuth.js or a hand-rolled OIDC client, or ' +
      'switch the auth axis to jwt-refresh / opaque-session.'
  );
}

export const adapter: AuthAdapter = {
  signIn: notImplemented,
  signOut: notImplemented,
  getSession: notImplemented
};

export default adapter;

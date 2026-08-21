import type { AuthAdapter } from '../types';
import { signIn, signOut, changePassword } from './auth';
import { getSession } from './session';

export const adapter: AuthAdapter = {
  signIn,
  signOut,
  getSession,
  changePassword
  // intentionally no refresh, no resetPassword for opaque-session
};

export default adapter;
export * from './types';

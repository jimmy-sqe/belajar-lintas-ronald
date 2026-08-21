import type { AuthAdapter } from '../types';
import { signIn, signOut, refresh } from './auth';
import { getSession } from './session';

export const adapter: AuthAdapter = {
  signIn,
  signOut,
  getSession,
  refresh
};

export default adapter;
export * from './types';

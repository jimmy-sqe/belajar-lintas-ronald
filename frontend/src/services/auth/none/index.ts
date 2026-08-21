import type { AuthAdapter, Session } from '../types';

export const adapter: AuthAdapter = {
  signIn: () => {
    throw new Error('Auth is disabled (auth=none). Configure an auth model in .boilerplate.yaml.');
  },
  signOut: () => Promise.resolve(),
  getSession: () => Promise.resolve<Session | null>(null)
};

export default adapter;

import { create } from 'zustand';

/**
 * Example Zustand store. Pattern reference for feature authors.
 *
 * Delete this file when adding real domain stores. Each domain typically
 * owns its own store file under src/store/<domain>.ts.
 */
interface ExampleState {
  count: number;
  increment: () => void;
  reset: () => void;
}

export const useExampleStore = create<ExampleState>((set) => ({
  count: 0,
  increment: () => set((s) => ({ count: s.count + 1 })),
  reset: () => set({ count: 0 })
}));

import { create } from 'zustand';
import { persist } from 'zustand/middleware';

const SYNC_STATUS = {
  IDLE: 'idle',
  SYNCING: 'syncing',
  ERROR: 'error',
  SUCCESS: 'success',
} as const;

type SyncStatus = (typeof SYNC_STATUS)[keyof typeof SYNC_STATUS];

interface PendingChange {
  id: string;
  type: 'sale' | 'product_update' | 'stock_adjustment';
  data: unknown;
  timestamp: number;
}

interface SyncState {
  isOnline: boolean;
  lastSync: number | null;
  pendingChanges: PendingChange[];
  status: SyncStatus;
}

interface SyncActions {
  setOnline: (isOnline: boolean) => void;
  setLastSync: (timestamp: number) => void;
  addPendingChange: (change: Omit<PendingChange, 'id' | 'timestamp'>) => void;
  removePendingChange: (id: string) => void;
  clearPendingChanges: () => void;
  setStatus: (status: SyncStatus) => void;
}

type SyncStore = SyncState & SyncActions;

const useSyncStore = create<SyncStore>()(
  persist(
    (set) => ({
      isOnline: navigator.onLine,
      lastSync: null,
      pendingChanges: [],
      status: SYNC_STATUS.IDLE,

      setOnline: (isOnline) => set({ isOnline }),

      setLastSync: (timestamp) => set({ lastSync: timestamp }),

      addPendingChange: (change) =>
        set((state) => ({
          pendingChanges: [
            ...state.pendingChanges,
            {
              ...change,
              id: crypto.randomUUID(),
              timestamp: Date.now(),
            },
          ],
        })),

      removePendingChange: (id) =>
        set((state) => ({
          pendingChanges: state.pendingChanges.filter((c) => c.id !== id),
        })),

      clearPendingChanges: () => set({ pendingChanges: [] }),

      setStatus: (status) => set({ status }),
    }),
    {
      name: 'sync-storage',
    }
  )
);

export { useSyncStore, SYNC_STATUS };
export type { SyncStatus, PendingChange, SyncState, SyncActions };

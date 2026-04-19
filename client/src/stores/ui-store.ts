import { create } from 'zustand';
import { persist } from 'zustand/middleware';

const THEME = {
  LIGHT: 'light',
  DARK: 'dark',
  SYSTEM: 'system',
} as const;

type Theme = (typeof THEME)[keyof typeof THEME];

interface UIState {
  sidebarOpen: boolean;
  theme: Theme;
}

interface UIActions {
  toggleSidebar: () => void;
  setSidebarOpen: (open: boolean) => void;
  setTheme: (theme: Theme) => void;
}

type UIStore = UIState & UIActions;

const useUIStore = create<UIStore>()(
  persist(
    (set) => ({
      sidebarOpen: true,
      theme: THEME.SYSTEM,

      toggleSidebar: () =>
        set((state) => ({ sidebarOpen: !state.sidebarOpen })),

      setSidebarOpen: (open) => set({ sidebarOpen: open }),

      setTheme: (theme) => set({ theme }),
    }),
    {
      name: 'ui-storage',
    }
  )
);

export { useUIStore, THEME };
export type { Theme, UIState, UIActions };

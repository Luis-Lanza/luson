import { create } from 'zustand';
import { persist } from 'zustand/middleware';

const ROLE = {
  ADMIN: 'admin',
  SELLER: 'seller',
  MANAGER: 'manager',
} as const;

type Role = (typeof ROLE)[keyof typeof ROLE];

interface User {
  id: string;
  email: string;
  name: string;
  role: Role;
  branchId: string | null;
}

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
}

interface AuthActions {
  login: (user: User, token: string) => void;
  logout: () => void;
  updateUser: (user: Partial<User>) => void;
}

type AuthStore = AuthState & AuthActions;

const useAuthStore = create<AuthStore>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,

      login: (user, token) =>
        set({
          user,
          token,
          isAuthenticated: true,
        }),

      logout: () =>
        set({
          user: null,
          token: null,
          isAuthenticated: false,
        }),

      updateUser: (userUpdate) =>
        set((state) => ({
          user: state.user ? { ...state.user, ...userUpdate } : null,
        })),
    }),
    {
      name: 'auth-storage',
    }
  )
);

export { useAuthStore, ROLE };
export type { User, Role, AuthState, AuthActions };

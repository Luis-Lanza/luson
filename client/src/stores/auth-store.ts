import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { User, UserRole, LoginResponse, RefreshResponse } from '../types';

export const ROLE = {
  ADMIN: 'admin' as const,
  ENCARGADO_ALMACEN: 'encargado_almacen' as const,
  CAJERO: 'cajero' as const,
} as const;

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

interface AuthState {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

interface AuthActions {
  login: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  refresh: () => Promise<void>;
  clearError: () => void;
  isAdmin: () => boolean;
  hasRole: (role: UserRole) => boolean;
  getRedirectPath: () => string;
}

type AuthStore = AuthState & AuthActions;

const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,

      login: async (username: string, password: string) => {
        set({ isLoading: true, error: null });

        try {
          const response = await fetch(`${API_URL}/api/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password }),
          });

          if (!response.ok) {
            const data = await response.json().catch(() => ({}));
            throw new Error(data.error || 'Credenciales inválidas');
          }

          const result: LoginResponse = await response.json();

          set({
            user: result.user,
            accessToken: result.access_token,
            refreshToken: result.refresh_token,
            isAuthenticated: true,
            isLoading: false,
            error: null,
          });
        } catch (error) {
          // For network errors (fetch throws TypeError on network failure)
          // or any error that's not a custom API error, show generic network message
          const isNetworkError = error instanceof TypeError || 
            (error instanceof Error && 
             error.message !== 'Credenciales inválidas' && 
             !error.message.includes('Unauthorized'));
          
          const message = isNetworkError 
            ? 'Error de conexión. Verifica tu red.'
            : (error instanceof Error ? error.message : 'Error desconocido');
          
          set({
            error: message,
            isLoading: false,
            isAuthenticated: false,
          });
        }
      },

      logout: async () => {
        const { accessToken } = get();

        // Try to call logout endpoint, but don't wait for it
        if (accessToken) {
          try {
            await fetch(`${API_URL}/api/auth/logout`, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${accessToken}`,
              },
            });
          } catch {
            // Ignore errors, we're logging out anyway
          }
        }

        // Clear all state
        set({
          user: null,
          accessToken: null,
          refreshToken: null,
          isAuthenticated: false,
          isLoading: false,
          error: null,
        });
      },

      refresh: async () => {
        const { refreshToken } = get();

        if (!refreshToken) {
          // No refresh token, logout
          set({
            user: null,
            accessToken: null,
            refreshToken: null,
            isAuthenticated: false,
          });
          return;
        }

        try {
          const response = await fetch(`${API_URL}/api/auth/refresh`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ refresh_token: refreshToken }),
          });

          if (!response.ok) {
            // Refresh failed, logout
            set({
              user: null,
              accessToken: null,
              refreshToken: null,
              isAuthenticated: false,
            });
            return;
          }

          const result: RefreshResponse = await response.json();

          set({
            accessToken: result.access_token,
          });
        } catch {
          // Refresh failed, logout
          set({
            user: null,
            accessToken: null,
            refreshToken: null,
            isAuthenticated: false,
          });
        }
      },

      clearError: () => {
        set({ error: null });
      },

      isAdmin: () => {
        const { user } = get();
        return user?.role === ROLE.ADMIN;
      },

      hasRole: (role: UserRole) => {
        const { user } = get();
        return user?.role === role;
      },

      getRedirectPath: () => {
        const { user } = get();
        
        if (!user) {
          return '/login';
        }

        if (user.role === ROLE.CAJERO) {
          return '/pos';
        }

        return '/dashboard';
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        refreshToken: state.refreshToken,
        user: state.user,
        // Don't persist accessToken - keep in memory only
        // Don't persist isAuthenticated - derive from accessToken presence
      }),
    }
  )
);

export { useAuthStore };
export type { AuthState, AuthActions };

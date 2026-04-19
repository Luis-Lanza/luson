import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { useAuthStore, ROLE } from './auth-store';

// Mock fetch globally
const mockFetch = vi.fn();
(globalThis as unknown as { fetch: typeof mockFetch }).fetch = mockFetch;

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
};
Object.defineProperty(window, 'localStorage', {
  value: localStorageMock,
  writable: true,
});

describe('auth-store', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Reset store state
    useAuthStore.setState({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,
    });
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('ROLE constants', () => {
    it('should have correct role values', () => {
      expect(ROLE.ADMIN).toBe('admin');
      expect(ROLE.ENCARGADO_ALMACEN).toBe('encargado_almacen');
      expect(ROLE.CAJERO).toBe('cajero');
    });
  });

  describe('initial state', () => {
    it('should have null user and tokens', () => {
      const state = useAuthStore.getState();
      expect(state.user).toBeNull();
      expect(state.accessToken).toBeNull();
      expect(state.refreshToken).toBeNull();
      expect(state.isAuthenticated).toBe(false);
      expect(state.isLoading).toBe(false);
      expect(state.error).toBeNull();
    });
  });

  describe('login', () => {
    it('should set loading state while logging in', async () => {
      mockFetch.mockImplementation(() => new Promise(() => {})); // Never resolves
      
      // Start login but don't await to check loading state
      useAuthStore.getState().login('testuser', 'password123');
      
      expect(useAuthStore.getState().isLoading).toBe(true);
      expect(useAuthStore.getState().error).toBeNull();
    });

    it('should store tokens and user on successful login', async () => {
      const mockUser = {
        id: 'user-1',
        username: 'testuser',
        role: 'cajero' as const,
        active: true,
        created_at: '2024-01-01T00:00:00Z',
      };
      
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          access_token: 'access-token-123',
          refresh_token: 'refresh-token-456',
          user: mockUser,
        }),
      });

      await useAuthStore.getState().login('testuser', 'password123');

      const state = useAuthStore.getState();
      expect(state.user).toEqual(mockUser);
      expect(state.accessToken).toBe('access-token-123');
      expect(state.refreshToken).toBe('refresh-token-456');
      expect(state.isAuthenticated).toBe(true);
      expect(state.isLoading).toBe(false);
      expect(state.error).toBeNull();
    });

    it('should call correct API endpoint with credentials', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          access_token: 'token',
          refresh_token: 'refresh',
          user: { id: '1', username: 'test', role: 'cajero', active: true, created_at: '' },
        }),
      });

      await useAuthStore.getState().login('myuser', 'mypassword');

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/auth/login'),
        expect.objectContaining({
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ username: 'myuser', password: 'mypassword' }),
        })
      );
    });

    it('should set error on failed login', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: async () => ({ error: 'Credenciales inválidas' }),
      });

      await useAuthStore.getState().login('wronguser', 'wrongpass');

      const state = useAuthStore.getState();
      expect(state.error).toBe('Credenciales inválidas');
      expect(state.isAuthenticated).toBe(false);
      expect(state.isLoading).toBe(false);
    });

    it('should handle network errors', async () => {
      // Network errors typically throw TypeError in fetch
      mockFetch.mockRejectedValueOnce(new TypeError('Failed to fetch'));

      await useAuthStore.getState().login('user', 'pass');

      const state = useAuthStore.getState();
      expect(state.error).toBe('Error de conexión. Verifica tu red.');
      expect(state.isAuthenticated).toBe(false);
    });
  });

  describe('logout', () => {
    it('should clear all state on logout', async () => {
      // First login
      useAuthStore.setState({
        user: { id: '1', username: 'test', role: 'admin', active: true, created_at: '' },
        accessToken: 'token',
        refreshToken: 'refresh',
        isAuthenticated: true,
      });

      mockFetch.mockResolvedValueOnce({ ok: true });

      await useAuthStore.getState().logout();

      const state = useAuthStore.getState();
      expect(state.user).toBeNull();
      expect(state.accessToken).toBeNull();
      expect(state.refreshToken).toBeNull();
      expect(state.isAuthenticated).toBe(false);
    });

    it('should call logout endpoint', async () => {
      useAuthStore.setState({
        accessToken: 'token',
        refreshToken: 'refresh',
      });

      mockFetch.mockResolvedValueOnce({ ok: true });

      await useAuthStore.getState().logout();

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/auth/logout'),
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({
            'Authorization': 'Bearer token',
          }),
        })
      );
    });

    it('should clear state even if logout endpoint fails', async () => {
      useAuthStore.setState({
        user: { id: '1', username: 'test', role: 'admin', active: true, created_at: '' },
        accessToken: 'token',
        refreshToken: 'refresh',
        isAuthenticated: true,
      });

      mockFetch.mockRejectedValueOnce(new Error('Network error'));

      await useAuthStore.getState().logout();

      const state = useAuthStore.getState();
      expect(state.isAuthenticated).toBe(false);
      expect(state.user).toBeNull();
    });
  });

  describe('refresh', () => {
    it('should refresh access token', async () => {
      useAuthStore.setState({
        refreshToken: 'old-refresh-token',
        accessToken: 'old-access-token',
      });

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          access_token: 'new-access-token',
        }),
      });

      await useAuthStore.getState().refresh();

      expect(useAuthStore.getState().accessToken).toBe('new-access-token');
      expect(useAuthStore.getState().refreshToken).toBe('old-refresh-token');
    });

    it('should call refresh endpoint with refresh token', async () => {
      useAuthStore.setState({
        refreshToken: 'my-refresh-token',
      });

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ access_token: 'new-token' }),
      });

      await useAuthStore.getState().refresh();

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/auth/refresh'),
        expect.objectContaining({
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ refresh_token: 'my-refresh-token' }),
        })
      );
    });

    it('should logout if refresh fails with 401', async () => {
      useAuthStore.setState({
        user: { id: '1', username: 'test', role: 'admin', active: true, created_at: '' },
        accessToken: 'token',
        refreshToken: 'expired-token',
        isAuthenticated: true,
      });

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: async () => ({ error: 'Token expirado' }),
      });

      await useAuthStore.getState().refresh();

      const state = useAuthStore.getState();
      expect(state.isAuthenticated).toBe(false);
      expect(state.user).toBeNull();
      expect(state.accessToken).toBeNull();
      expect(state.refreshToken).toBeNull();
    });
  });

  describe('clearError', () => {
    it('should clear error state', () => {
      useAuthStore.setState({ error: 'Some error' });
      
      useAuthStore.getState().clearError();
      
      expect(useAuthStore.getState().error).toBeNull();
    });
  });

  describe('isAdmin selector', () => {
    it('should return true for admin user', () => {
      useAuthStore.setState({
        user: { id: '1', username: 'admin', role: 'admin', active: true, created_at: '' },
      });

      expect(useAuthStore.getState().isAdmin()).toBe(true);
    });

    it('should return false for non-admin user', () => {
      useAuthStore.setState({
        user: { id: '1', username: 'cashier', role: 'cajero', active: true, created_at: '' },
      });

      expect(useAuthStore.getState().isAdmin()).toBe(false);
    });

    it('should return false when no user', () => {
      useAuthStore.setState({ user: null });
      expect(useAuthStore.getState().isAdmin()).toBe(false);
    });
  });

  describe('hasRole selector', () => {
    it('should return true when user has specified role', () => {
      useAuthStore.setState({
        user: { id: '1', username: 'cashier', role: 'cajero', active: true, created_at: '' },
      });

      expect(useAuthStore.getState().hasRole('cajero')).toBe(true);
      expect(useAuthStore.getState().hasRole('admin')).toBe(false);
    });
  });

  describe('getRedirectPath selector', () => {
    it('should return /pos for cajero', () => {
      useAuthStore.setState({
        user: { id: '1', username: 'cashier', role: 'cajero', active: true, created_at: '' },
      });
      expect(useAuthStore.getState().getRedirectPath()).toBe('/pos');
    });

    it('should return /dashboard for admin', () => {
      useAuthStore.setState({
        user: { id: '1', username: 'admin', role: 'admin', active: true, created_at: '' },
      });
      expect(useAuthStore.getState().getRedirectPath()).toBe('/dashboard');
    });

    it('should return /dashboard for encargado_almacen', () => {
      useAuthStore.setState({
        user: { id: '1', username: 'manager', role: 'encargado_almacen', active: true, created_at: '' },
      });
      expect(useAuthStore.getState().getRedirectPath()).toBe('/dashboard');
    });

    it('should return /login when no user', () => {
      useAuthStore.setState({ user: null });
      expect(useAuthStore.getState().getRedirectPath()).toBe('/login');
    });
  });
});

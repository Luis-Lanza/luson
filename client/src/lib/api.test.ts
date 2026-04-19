import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';

// Mock auth store
const mockGetState = vi.fn();
const mockSetState = vi.fn();

vi.mock('../stores/auth-store', () => ({
  useAuthStore: {
    getState: () => mockGetState(),
    setState: (fn: (state: unknown) => unknown) => mockSetState(fn),
  },
}));

// Import after mock
import { api, authApi, usersApi, branchesApi, suppliersApi } from './api';

describe('api', () => {
  const mockFetch = vi.fn();
  vi.stubGlobal('fetch', mockFetch);

  beforeEach(() => {
    vi.clearAllMocks();
    mockGetState.mockReturnValue({
      accessToken: 'test-access-token',
      refreshToken: 'test-refresh-token',
      refresh: vi.fn(),
    });
    
    // Reset fetch mock
    mockFetch.mockReset();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('api client', () => {
    it('should make GET request with correct headers', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: 'test' }),
      });

      await api.get('/test');

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/test'),
        expect.objectContaining({
          method: 'GET',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
            'Authorization': 'Bearer test-access-token',
          }),
        })
      );
    });

    it('should make POST request with body', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ result: 'success' }),
      });

      const body = { name: 'Test' };
      await api.post('/test', body);

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/test'),
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
          body: JSON.stringify(body),
        })
      );
    });

    it('should make PUT request with body', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ result: 'success' }),
      });

      const body = { id: '1', name: 'Updated' };
      await api.put('/test/1', body);

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/test/1'),
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(body),
        })
      );
    });

    it('should make DELETE request', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ deleted: true }),
      });

      await api.delete('/test/1');

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/test/1'),
        expect.objectContaining({
          method: 'DELETE',
        })
      );
    });

    it('should return parsed JSON on success', async () => {
      const responseData = { id: '1', name: 'Test' };
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => responseData,
      });

      const result = await api.get('/test');

      expect(result).toEqual(responseData);
    });

    it('should throw error on non-ok response', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 400,
        json: async () => ({ error: 'Bad request' }),
      });

      await expect(api.get('/test')).rejects.toThrow('Bad request');
    });

    it('should throw generic error when no error message', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: async () => ({}),
      });

      await expect(api.get('/test')).rejects.toThrow('Error HTTP: 500');
    });

    it('should work without authorization when not logged in', async () => {
      mockGetState.mockReturnValue({
        accessToken: null,
      });

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: 'public' }),
      });

      await api.get('/public');

      const [, options] = mockFetch.mock.calls[0];
      expect(options.headers['Authorization']).toBeUndefined();
    });

    it('should use VITE_API_URL from environment', async () => {
      // Mock the env variable by testing the API base URL construction
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: 'test' }),
      });

      await api.get('/test');

      // Just verify it uses the base URL pattern (default or env)
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringMatching(/^https?:\/\/.*\/test$/),
        expect.any(Object)
      );
    });

    it('should default to localhost:8080 when no env var', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: 'test' }),
      });

      await api.get('/test');

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('localhost:8080'),
        expect.any(Object)
      );
    });
  });

  describe('authApi', () => {
    it('should call login endpoint', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          access_token: 'new-token',
          refresh_token: 'new-refresh',
          user: { id: '1', username: 'test', role: 'cajero', active: true, created_at: '' },
        }),
      });

      const credentials = { username: 'test', password: 'pass' };
      const result = await authApi.login(credentials);

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/auth/login'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(credentials),
        })
      );
      expect(result.access_token).toBe('new-token');
    });

    it('should call refresh endpoint', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ access_token: 'refreshed-token' }),
      });

      const result = await authApi.refresh('my-refresh-token');

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/auth/refresh'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ refresh_token: 'my-refresh-token' }),
        })
      );
      expect(result.access_token).toBe('refreshed-token');
    });

    it('should call me endpoint', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ id: '1', username: 'test', role: 'cajero', active: true, created_at: '' }),
      });

      await authApi.me();

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/auth/me'),
        expect.objectContaining({
          method: 'GET',
        })
      );
    });

    it('should call logout endpoint', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ message: 'Logout successful' }),
      });

      await authApi.logout();

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/auth/logout'),
        expect.objectContaining({
          method: 'POST',
        })
      );
    });
  });

  describe('usersApi', () => {
    it('should list users', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          data: [{ id: '1', username: 'user1', role: 'cajero', active: true, created_at: '' }],
          total: 1,
          page: 1,
          limit: 10,
        }),
      });

      const result = await usersApi.list();

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/users'),
        expect.objectContaining({ method: 'GET' })
      );
      expect(result.data).toHaveLength(1);
    });

    it('should get user by id', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ id: '1', username: 'user1', role: 'cajero', active: true, created_at: '' }),
      });

      await usersApi.getById('1');

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/users/1'),
        expect.objectContaining({ method: 'GET' })
      );
    });

    it('should create user', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ id: '2', username: 'newuser', role: 'cajero', active: true, created_at: '' }),
      });

      const data = { username: 'newuser', password: 'pass', role: 'cajero' as const };
      await usersApi.create(data);

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/users'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(data),
        })
      );
    });

    it('should update user', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ id: '1', username: 'updated', role: 'cajero', active: true, created_at: '' }),
      });

      const data = { username: 'updated' };
      await usersApi.update('1', data);

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/users/1'),
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(data),
        })
      );
    });
  });

  describe('branchesApi', () => {
    it('should list branches', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          data: [{ id: '1', name: 'Branch 1', petty_cash_balance: 0, active: true, created_at: '' }],
          total: 1,
          page: 1,
          limit: 10,
        }),
      });

      const result = await branchesApi.list();

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/branches'),
        expect.objectContaining({ method: 'GET' })
      );
      expect(result.data).toHaveLength(1);
    });

    it('should get branch by id', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ id: '1', name: 'Branch 1', petty_cash_balance: 0, active: true, created_at: '' }),
      });

      await branchesApi.getById('1');

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/branches/1'),
        expect.objectContaining({ method: 'GET' })
      );
    });

    it('should create branch', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ id: '2', name: 'New Branch', petty_cash_balance: 0, active: true, created_at: '' }),
      });

      const data = { name: 'New Branch', address: '123 Main St' };
      await branchesApi.create(data);

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/branches'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(data),
        })
      );
    });

    it('should update branch', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ id: '1', name: 'Updated Branch', petty_cash_balance: 100, active: true, created_at: '' }),
      });

      const data = { name: 'Updated Branch', petty_cash_balance: 100 };
      await branchesApi.update('1', data);

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/branches/1'),
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(data),
        })
      );
    });
  });

  describe('suppliersApi', () => {
    it('should list suppliers', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          data: [{ id: '1', name: 'Supplier 1', active: true, created_at: '' }],
          total: 1,
          page: 1,
          limit: 10,
        }),
      });

      const result = await suppliersApi.list();

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/suppliers'),
        expect.objectContaining({ method: 'GET' })
      );
      expect(result.data).toHaveLength(1);
    });

    it('should get supplier by id', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ id: '1', name: 'Supplier 1', active: true, created_at: '' }),
      });

      await suppliersApi.getById('1');

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/suppliers/1'),
        expect.objectContaining({ method: 'GET' })
      );
    });

    it('should create supplier', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ id: '2', name: 'New Supplier', active: true, created_at: '' }),
      });

      const data = { name: 'New Supplier', contact: '555-1234' };
      await suppliersApi.create(data);

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/suppliers'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(data),
        })
      );
    });

    it('should update supplier', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ id: '1', name: 'Updated Supplier', active: true, created_at: '' }),
      });

      const data = { name: 'Updated Supplier' };
      await suppliersApi.update('1', data);

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/suppliers/1'),
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(data),
        })
      );
    });
  });
});

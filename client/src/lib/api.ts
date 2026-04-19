import { useAuthStore } from '../stores/auth-store';
import type {
  User,
  Branch,
  Supplier,
  Product,
  Stock,
  PurchaseBatch,
  Transfer,
  LoginCredentials,
  LoginResponse,
  RefreshResponse,
  CreateUserRequest,
  UpdateUserRequest,
  CreateBranchRequest,
  UpdateBranchRequest,
  CreateSupplierRequest,
  UpdateSupplierRequest,
  PaginatedResponse,
} from '../types';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

// Track pending requests for queueing during token refresh
let isRefreshing = false;
let refreshQueue: Array<() => void> = [];

async function processQueue(error: Error | null = null) {
  refreshQueue.forEach((callback) => {
    if (error) {
      // If there was an error, we could reject the promises
      // For now, just call the callbacks
    }
    callback();
  });
  refreshQueue = [];
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const authState = useAuthStore.getState();
  const accessToken = authState.accessToken;

  const url = `${API_URL}${path}`;
  
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...((options.headers as Record<string, string>) || {}),
  };

  if (accessToken) {
    headers['Authorization'] = `Bearer ${accessToken}`;
  }

  try {
    const response = await fetch(url, {
      ...options,
      headers,
    });

    // Handle 401 - Token expired, try to refresh
    if (response.status === 401 && accessToken) {
      // If we're already refreshing, queue this request
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          refreshQueue.push(() => {
            request<T>(path, options).then(resolve).catch(reject);
          });
        });
      }

      isRefreshing = true;

      try {
        await authState.refresh();
        const newToken = useAuthStore.getState().accessToken;
        
        if (!newToken) {
          throw new Error('Session expired');
        }

        // Retry the original request with new token
        const retryHeaders = {
          ...headers,
          'Authorization': `Bearer ${newToken}`,
        };

        const retryResponse = await fetch(url, {
          ...options,
          headers: retryHeaders,
        });

        if (!retryResponse.ok) {
          const errorData = await retryResponse.json().catch(() => ({}));
          throw new Error(errorData.error || `Error HTTP: ${retryResponse.status}`);
        }

        const data = await retryResponse.json();
        processQueue();
        return data as T;
      } catch (refreshError) {
        processQueue(refreshError as Error);
        throw refreshError;
      } finally {
        isRefreshing = false;
      }
    }

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || `Error HTTP: ${response.status}`);
    }

    return await response.json() as T;
  } catch (error) {
    if (error instanceof Error) {
      throw error;
    }
    throw new Error('Error de red');
  }
}

// Generic API methods
export const api = {
  get: <T>(path: string) => request<T>(path, { method: 'GET' }),
  post: <T>(path: string, body?: unknown) =>
    request<T>(path, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    }),
  put: <T>(path: string, body?: unknown) =>
    request<T>(path, {
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined,
    }),
  delete: <T>(path: string) => request<T>(path, { method: 'DELETE' }),
};

// Auth API
export const authApi = {
  login: (credentials: LoginCredentials) =>
    api.post<LoginResponse>('/api/auth/login', credentials),
  refresh: (refreshToken: string) =>
    api.post<RefreshResponse>('/api/auth/refresh', { refresh_token: refreshToken }),
  me: () => api.get<User>('/api/auth/me'),
  logout: () => api.post<void>('/api/auth/logout'),
};

// Users API
export const usersApi = {
  list: (params?: { page?: number; limit?: number }) => {
    const queryParams = params
      ? '?' + new URLSearchParams(params as Record<string, string>).toString()
      : '';
    return api.get<PaginatedResponse<User>>(`/api/users${queryParams}`);
  },
  getById: (id: string) => api.get<User>(`/api/users/${id}`),
  create: (data: CreateUserRequest) => api.post<User>('/api/users', data),
  update: (id: string, data: UpdateUserRequest) => api.put<User>(`/api/users/${id}`, data),
};

// Branches API
export const branchesApi = {
  list: (params?: { page?: number; limit?: number }) => {
    const queryParams = params
      ? '?' + new URLSearchParams(params as Record<string, string>).toString()
      : '';
    return api.get<PaginatedResponse<Branch>>(`/api/branches${queryParams}`);
  },
  getById: (id: string) => api.get<Branch>(`/api/branches/${id}`),
  create: (data: CreateBranchRequest) => api.post<Branch>('/api/branches', data),
  update: (id: string, data: UpdateBranchRequest) => api.put<Branch>(`/api/branches/${id}`, data),
};

// Suppliers API
export const suppliersApi = {
  list: (params?: { page?: number; limit?: number }) => {
    const queryParams = params
      ? '?' + new URLSearchParams(params as Record<string, string>).toString()
      : '';
    return api.get<PaginatedResponse<Supplier>>(`/api/suppliers${queryParams}`);
  },
  getById: (id: string) => api.get<Supplier>(`/api/suppliers/${id}`),
  create: (data: CreateSupplierRequest) => api.post<Supplier>('/api/suppliers', data),
  update: (id: string, data: UpdateSupplierRequest) => api.put<Supplier>(`/api/suppliers/${id}`, data),
};

// Products API
export const productsApi = {
  list: (params?: Record<string, string>) =>
    api.get<Product[]>('/api/products' + (params ? '?' + new URLSearchParams(params) : '')),
  getById: (id: string) => api.get<Product>(`/api/products/${id}`),
  create: (data: Omit<Product, 'id' | 'created_at'>) => api.post<Product>('/api/products', data),
  update: (id: string, data: Partial<Product>) => api.put<Product>(`/api/products/${id}`, data),
};

// Stock API
export const stockApi = {
  list: (params?: Record<string, string>) =>
    api.get<Stock[]>('/api/stock' + (params ? '?' + new URLSearchParams(params) : '')),
  getLowAlerts: () => api.get<Stock[]>('/api/stock/low'),
  adjust: (data: { product_id: string; product_type: string; location_type: string; location_id: string; quantity: number; reason: string }) =>
    api.post('/api/stock/adjust', data),
};

// Purchase Batches API
export const purchaseBatchesApi = {
  list: () => api.get<PurchaseBatch[]>('/api/purchase-batches'),
  getById: (id: string) => api.get<PurchaseBatch>(`/api/purchase-batches/${id}`),
  create: (data: { supplier_id: string; purchase_date: string; notes?: string; details: { product_id: string; product_type: string; quantity: number; unit_cost: number }[] }) =>
    api.post<PurchaseBatch>('/api/purchase-batches', data),
  process: (id: string) => api.post<void>(`/api/purchase-batches/${id}/process`),
};

// Transfers API
export const transfersApi = {
  list: (params?: Record<string, string>) =>
    api.get<Transfer[]>('/api/transfers' + (params ? '?' + new URLSearchParams(params) : '')),
  getById: (id: string) => api.get<Transfer>(`/api/transfers/${id}`),
  create: (data: { origin_type: string; origin_id: string; destination_type: string; destination_id: string; transfer_type: string; details: { product_id?: string; product_type?: string; quantity: number; liters?: number }[] }) =>
    api.post<Transfer>('/api/transfers', data),
  approve: (id: string) => api.post<void>(`/api/transfers/${id}/approve`),
  reject: (id: string, reason: string) => api.post<void>(`/api/transfers/${id}/reject`, { reason }),
  ship: (id: string) => api.post<void>(`/api/transfers/${id}/ship`),
  receive: (id: string) => api.post<void>(`/api/transfers/${id}/receive`),
  push: (data: { destination_type: string; destination_id: string; details: { product_id?: string; product_type?: string; quantity: number; liters?: number }[] }) =>
    api.post<Transfer>('/api/transfers/push', data),
};

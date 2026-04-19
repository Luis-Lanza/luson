export type UserRole = 'admin' | 'encargado_almacen' | 'cajero';

export interface User {
  id: string;
  username: string;
  role: UserRole;
  branch_id?: string;
  active: boolean;
  created_at: string;
}

export interface Branch {
  id: string;
  name: string;
  address?: string;
  petty_cash_balance: number;
  active: boolean;
  created_at: string;
}

export interface Supplier {
  id: string;
  name: string;
  contact?: string;
  address?: string;
  active: boolean;
  created_at: string;
}

export interface LoginCredentials {
  username: string;
  password: string;
}

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  user: User;
}

export interface RefreshResponse {
  access_token: string;
}

export interface ApiResponse<T> {
  data?: T;
  error?: string;
  message?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
}

export interface CreateUserRequest {
  username: string;
  password: string;
  role: UserRole;
  branch_id?: string;
}

export interface UpdateUserRequest {
  username?: string;
  password?: string;
  role?: UserRole;
  branch_id?: string;
  active?: boolean;
}

export interface CreateBranchRequest {
  name: string;
  address?: string;
  petty_cash_balance?: number;
}

export interface UpdateBranchRequest {
  name?: string;
  address?: string;
  petty_cash_balance?: number;
  active?: boolean;
}

export interface CreateSupplierRequest {
  name: string;
  contact?: string;
  address?: string;
}

export interface UpdateSupplierRequest {
  name?: string;
  contact?: string;
  address?: string;
  active?: boolean;
}

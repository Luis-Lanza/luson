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

// Product Types
export type ProductType = 'bateria' | 'accesorio';
export type BatteryType = 'seca' | 'liquida';
export type Polarity = 'izquierda' | 'derecha';
export type VehicleType = 'auto' | 'moto' | 'otro';
export type TransferStatus = 'pendiente' | 'aprobada' | 'rechazada' | 'enviada' | 'recibida';

export interface Product {
  id: string;
  name: string;
  description?: string;
  product_type: ProductType;
  brand?: string;
  model?: string;
  voltage?: number;
  amperage?: number;
  battery_type?: BatteryType;
  polarity?: Polarity;
  acid_liters?: number;
  vehicle_type?: VehicleType;
  min_sale_price: number;
  active: boolean;
  created_at: string;
}

export interface Stock {
  id: string;
  product_id: string;
  product_type: ProductType;
  location_type: string;
  location_id: string;
  quantity: number;
  min_stock_alert: number;
  updated_at: string;
}

export interface PurchaseBatch {
  id: string;
  supplier_id: string;
  purchase_date: string;
  notes?: string;
  total_cost: number;
  processed: boolean;
  details?: PurchaseBatchDetail[];
  created_at: string;
}

export interface PurchaseBatchDetail {
  product_id: string;
  product_type: ProductType;
  quantity: number;
  unit_cost: number;
}

export interface Transfer {
  id: string;
  origin_type: string;
  origin_id: string;
  destination_type: string;
  destination_id: string;
  status: TransferStatus;
  rejection_reason?: string;
  transfer_type: string;
  details?: TransferDetail[];
  created_at: string;
}

export interface TransferDetail {
  product_id?: string;
  product_type?: ProductType;
  quantity: number;
  liters?: number;
}

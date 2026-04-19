import { z, type ZodType } from 'zod';

// Base schemas
export const UserSchema = z.object({
  id: z.string(),
  email: z.string().email(),
  name: z.string(),
  role: z.enum(['admin', 'seller', 'manager']),
  branchId: z.string().nullable(),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
});

export const BranchSchema = z.object({
  id: z.string(),
  name: z.string(),
  address: z.string(),
  phone: z.string().optional(),
  isActive: z.boolean(),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
});

export const ProductSchema = z.object({
  id: z.string(),
  sku: z.string(),
  name: z.string(),
  description: z.string().optional(),
  price: z.number().positive(),
  cost: z.number().positive(),
  stock: z.number().int().min(0),
  minStock: z.number().int().min(0),
  category: z.string(),
  isActive: z.boolean(),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
});

export const SaleItemSchema = z.object({
  id: z.string(),
  productId: z.string(),
  productName: z.string(),
  quantity: z.number().int().positive(),
  unitPrice: z.number().positive(),
  total: z.number().positive(),
});

export const SaleSchema = z.object({
  id: z.string(),
  items: z.array(SaleItemSchema),
  subtotal: z.number().positive(),
  tax: z.number().min(0),
  total: z.number().positive(),
  paymentMethod: z.enum(['cash', 'card', 'transfer']),
  customerName: z.string().optional(),
  customerPhone: z.string().optional(),
  sellerId: z.string(),
  branchId: z.string(),
  createdAt: z.string().datetime(),
  syncedAt: z.string().datetime().optional(),
});

// API response schemas
export function createApiResponseSchema<T extends ZodType>(dataSchema: T) {
  return z.object({
    success: z.boolean(),
    data: dataSchema,
    message: z.string().optional(),
  });
}

export const ApiErrorSchema = z.object({
  success: z.literal(false),
  error: z.object({
    code: z.string(),
    message: z.string(),
    details: z.record(z.string(), z.unknown()).optional(),
  }),
});

// Type inference
export type User = z.infer<typeof UserSchema>;
export type Branch = z.infer<typeof BranchSchema>;
export type Product = z.infer<typeof ProductSchema>;
export type SaleItem = z.infer<typeof SaleItemSchema>;
export type Sale = z.infer<typeof SaleSchema>;
export type ApiError = z.infer<typeof ApiErrorSchema>;

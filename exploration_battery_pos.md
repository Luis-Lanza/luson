# Exploration: Battery POS Phase 2 - Productos + Stock

## 1. Recommended Entities and Their Fields

### Product Entity (Unified Table Approach)
Based on PRD sections 4.3 and 4.4, I recommend a single `products` table with type discrimination rather than separate tables for batteries and accessories. This approach:
- Simplifies queries across product types
- Reduces joins in stock and sale operations
- Follows existing patterns in the codebase
- Makes it easier to add new product types in the future

**Product Fields:**
- id: UUID (primary key)
- name: string (for accessories) OR model: string + brand: string (for batteries) - I'll use name as generic field
- description: string? (optional, mainly for accessories)
- product_type: enum (`bateria`, `accesorio`) - discrimination field
- battery_specific fields (nullable when product_type = 'accesorio'):
  - brand: string?
  - model: string?
  - voltage: decimal?
  - amperage: decimal?
  - battery_type: enum? (`seca`, `liquida`)
  - polarity: enum? (`izquierda`, `derecha`)
  - acid_liters: decimal?
  - vehicle_type: enum? (`auto`, `moto`, `otro`)
- min_sale_price: decimal
- effective_date: timestamp? (for price changes)
- previous_price: decimal? (when effective_date is in future)
- active: bool
- created_at: timestamp
- created_by: UUID (user who created)

### Stock Entity
**Stock Fields:**
- id: UUID (primary key)
- product_id: UUID (foreign key to products)
- product_type: enum (`bateria`, `accesorio`) - denormalized for performance
- location_type: enum (`almacen`, `sucursal`)
- location_id: UUID (references branches.id)
- quantity: int
- min_stock_alert: int
- updated_at: timestamp

### Purchase Batch Entity (Lote de Compra)
**PurchaseBatch Fields:**
- id: UUID (primary key)
- supplier_id: UUID (foreign key to suppliers)
- purchase_date: timestamp
- notes: string?
- created_by: UUID
- created_at: timestamp

### Purchase Batch Detail Entity (Detalle de Lote)
**PurchaseBatchDetail Fields:**
- id: UUID (primary key)
- batch_id: UUID (foreign key to purchase_batches)
- product_id: UUID (foreign key to products)
- product_type: enum (`bateria`, `accesorio`) - denormalized
- quantity: int
- unit_cost: decimal (admin-only visibility)
- created_at: timestamp

### Transfer Entity
**Transfer Fields:**
- id: UUID (primary key)
- origin_type: enum (`almacen`, `sucursal`)
- origin_id: UUID (references branches.id)
- destination_type: enum (`almacen`, `sucursal`)
- destination_id: UUID (references branches.id)
- status: enum (`pendiente`, `aprobada`, `rechazada`, `enviada`, `recibida`)
- rejection_reason: string?
- transfer_type: enum (`stock`, `chatarra`, `bateria_fallada`, `acido`) - focusing on stock for Phase 2
- created_by: UUID
- created_at: timestamp
- approved_by: UUID?
- approved_at: timestamp?
- shipped_at: timestamp?
- received_at: timestamp?
- received_by: UUID?

### Transfer Detail Entity (for stock transfers)
**TransferDetail Fields:**
- id: UUID (primary key)
- transfer_id: UUID (foreign key to transfers)
- product_id: UUID? (nullable for non-stock transfers)
- product_type: enum? (`bateria`, `accesorio`)
- quantity: decimal
- created_at: timestamp

## 2. API Endpoints

### Product Management
- GET /api/products - List products (with filters: product_type, active, search)
- GET /api/products/{id} - Get product by ID
- POST /api/products - Create product (battery or accessory)
- PUT /api/products/{id} - Update product
- DELETE /api/products/{id} - Soft delete (set active=false)
- GET /api/products/{id}/stock - Get stock levels across all locations

### Stock Management
- GET /api/stock - List stock (with filters: location_type, location_id, product_id, low_stock)
- GET /api/stock/{id} - Get stock item by ID
- POST /api/stock - Adjust stock manually (for inicialización)
- PUT /api/stock/{id} - Update stock levels (typically done by system)
- GET /api/stock/low-alerts - Get products with stock below threshold

### Purchase Batch Management
- GET /api/purchase-batches - List purchase batches (with filters: supplier, date range)
- GET /api/purchase-batches/{id} - Get purchase batch with details
- POST /api/purchase-batches - Create purchase batch with details
- PUT /api/purchase-batches/{id} - Update purchase batch (before processing)
- POST /api/purchase-batches/{id}/process - Process batch (add to stock)
- DELETE /api/purchase-batches/{id} - Cancel batch (if not processed)

### Transfer Management
- GET /api/transfers - List transfers (with filters: status, origin/destination, type)
- GET /api/transfers/{id} - Get transfer with details
- POST /api/transfers - Create transfer request (branch to branch/branch to almacen/almacen to branch)
- PUT /api/transfers/{id}/approve - Approve transfer (almacén or destino)
- PUT /api/transfers/{id}/reject - Reject transfer (with reason)
- PUT /api/transfers/{id}/ship - Mark as shipped
- PUT /api/transfers/{id}/receive - Mark as received
- POST /api/transfers/push - Almacén push stock to sucursal (direct transfer)

## 3. Database Schema Changes

### New Tables

```sql
-- Products table (unified for batteries and accessories)
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    product_type VARCHAR(50) NOT NULL CHECK (product_type IN ('bateria', 'accesorio')),
    -- Battery-specific fields (nullable for accessories)
    brand VARCHAR(100),
    model VARCHAR(100),
    voltage DECIMAL(5,2),
    amperage DECIMAL(5,2),
    battery_type VARCHAR(20) CHECK (battery_type IN ('seca', 'liquida')),
    polarity VARCHAR(20) CHECK (polarity IN ('izquierda', 'derecha')),
    acid_liters DECIMAL(4,2),
    vehicle_type VARCHAR(20) CHECK (vehicle_type IN ('auto', 'moto', 'otro')),
    -- Pricing fields
    min_sale_price DECIMAL(12,2) NOT NULL,
    effective_date TIMESTAMP WITH TIME ZONE,
    previous_price DECIMAL(12,2),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id)
);

-- Indexes for products
CREATE INDEX IF NOT EXISTS idx_products_type_active ON products(product_type, active);
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
CREATE INDEX IF NOT EXISTS idx_products_battery_fields ON products(brand, model) WHERE product_type = 'bateria';

-- Stock table
CREATE TABLE IF NOT EXISTS stock (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    product_type VARCHAR(50) NOT NULL CHECK (product_type IN ('bateria', 'accesorio')),
    location_type VARCHAR(20) NOT NULL CHECK (location_type IN ('almacen', 'sucursal')),
    location_id UUID NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 0,
    min_stock_alert INTEGER NOT NULL DEFAULT 5,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(product_id, location_type, location_id)
);

-- Indexes for stock
CREATE INDEX IF NOT EXISTS idx_stock_location ON stock(location_type, location_id);
CREATE INDEX IF NOT EXISTS idx_stock_product ON stock(product_id);
CREATE INDEX IF NOT EXISTS idx_stock_low_alert ON stock(quantity, min_stock_alert) WHERE quantity < min_stock_alert;

-- Purchase batches table
CREATE TABLE IF NOT EXISTS purchase_batches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    supplier_id UUID NOT NULL REFERENCES suppliers(id),
    purchase_date TIMESTAMP WITH TIME ZONE NOT NULL,
    notes TEXT,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for purchase batches
CREATE INDEX IF NOT EXISTS idx_purchase_batches_supplier ON purchase_batches(supplier_id);
CREATE INDEX IF NOT EXISTS idx_purchase_batches_date ON purchase_batches(purchase_date);

-- Purchase batch details table
CREATE TABLE IF NOT EXISTS purchase_batch_details (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    batch_id UUID NOT NULL REFERENCES purchase_batches(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    product_type VARCHAR(50) NOT NULL CHECK (product_type IN ('bateria', 'accesorio')),
    quantity INTEGER NOT NULL,
    unit_cost DECIMAL(12,4) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for batch details
CREATE INDEX IF NOT EXISTS idx_batch_details_batch ON purchase_batch_details(batch_id);
CREATE INDEX IF NOT EXISTS idx_batch_details_product ON purchase_batch_details(product_id);

-- Transfers table
CREATE TABLE IF NOT EXISTS transfers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    origin_type VARCHAR(20) NOT NULL CHECK (origin_type IN ('almacen', 'sucursal')),
    origin_id UUID NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    destination_type VARCHAR(20) NOT NULL CHECK (destination_type IN ('almacen', 'sucursal')),
    destination_id UUID NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pendiente', 'aprobada', 'rechazada', 'enviada', 'recibida')),
    rejection_reason TEXT,
    transfer_type VARCHAR(20) NOT NULL CHECK (transfer_type IN ('stock', 'chatarra', 'bateria_fallada', 'acido')),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP WITH TIME ZONE,
    shipped_at TIMESTAMP WITH TIME ZONE,
    received_at TIMESTAMP WITH TIME ZONE,
    received_by UUID REFERENCES users(id)
);

-- Indexes for transfers
CREATE INDEX IF NOT EXISTS idx_transfers_status ON transfers(status);
CREATE INDEX IF NOT EXISTS idx_transfers_origin ON transfers(origin_type, origin_id);
CREATE INDEX IF NOT EXISTS idx_transfers_destination ON transfers(destination_type, destination_id);
CREATE INDEX IF NOT EXISTS idx_transfers_type ON transfers(transfer_type);

-- Transfer details table
CREATE TABLE IF NOT EXISTS transfer_details (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transfer_id UUID NOT NULL REFERENCES transfers(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id) ON DELETE SET NULL,
    product_type VARCHAR(50) CHECK (product_type IN ('bateria', 'accesorio')),
    quantity DECIMAL(12,2) NOT NULL,
    liters DECIMAL(8,3), -- for acid transfers
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for transfer details
CREATE INDEX IF NOT EXISTS idx_transfer_details_transfer ON transfer_details(transfer_id);
CREATE INDEX IF NOT EXISTS idx_transfer_details_product ON transfer_details(product_id);
```

### Modified Tables
Add foreign key constraints to existing tables if needed, but based on current schema, no modifications needed to existing tables.

## 4. Test Strategy

### Backend Testing (TDD Approach)
Following STRICT TDD: write tests FIRST, then implement.

#### Unit Tests
- Domain entity validation (Product.IsValid(), Stock.IsValid(), etc.)
- Service layer tests (ProductService, StockService, PurchaseBatchService, TransferService)
- Repository tests (using testcontainers or in-memory SQLite for Postgres)

#### Integration Tests
- API endpoint tests (using httptest)
- Full flow tests: Create product → Check stock → Create purchase batch → Process batch → Verify stock increase
- Transfer flow tests: Request transfer → Approve → Ship → Receive → Verify stock movement

#### Test Structure
```
/internal/domain/product_test.go
/internal/domain/stock_test.go
/internal/application/product_service_test.go
/internal/application/stock_service_test.go
/internal/application/purchase_batch_service_test.go
/internal/application/transfer_service_test.go
/internal/infrastructure/postgres/product_repo_test.go
/internal/infrastructure/postgres/stock_repo_test.go
/internal/infrastructure/http/handlers/product_handler_test.go
/internal/infrastructure/http/handlers/stock_handler_test.go
```

### Frontend Testing
- Unit tests for product/form components
- Integration tests for product CRUD pages
- E2E tests for critical flows (using Playwright or similar)
- Test stock alerts and notifications

## 5. Risks and Gotchas

### Technical Risks
1. **Data Migration**: No existing product data to migrate, but need to ensure backward compatibility
2. **Enum Management**: Need to handle enum values consistently across Go, SQL, and TypeScript
3. **Performance**: Stock queries need to be optimized with proper indexes
4. **Concurrency**: Stock updates need to handle race conditions (use SELECT FOR UPDATE or similar)
5. **Offline Sync**: Frontend needs to handle syncing of product/stock data when offline

### Business Logic Gotchas
1. **Price Effective Dates**: Need to handle price changes where effective_date is in the future
2. **Battery Validation**: Liquid batteries require acid_liters, seca batteries should not have acid_liters
3. **Stock Deduction**: Ensure stock never goes negative (except for acid which allows negative)
4. **Transfer Validation**: Cannot transfer to same location, need valid quantities
5. **Batch Processing**: Unit cost should be stored but not exposed to cashiers (admin-only)

### Implementation Challenges
1. **Type Discrimination**: Handling product-type-specific fields in a unified table
2. **API Design**: Balancing RESTfulness with complex nested resources (batch details, transfer details)
3. **Frontend State Management**: Managing product forms with conditional fields based on type
4. **Real-time Updates**: Stock changes need to reflect in real-time across locations when online
5. **Notification Triggers**: Low stock alerts need to trigger appropriately without spamming

### Recommendations
1. Start with domain entities and validation logic
2. Implement repository interfaces and Postgres implementations
3. Build service layer with business rules
4. Create API handlers
5. Develop frontend components last
6. Use existing patterns in the codebase (look at user/branch/supplier implementations)
7. Follow the same testing patterns already established
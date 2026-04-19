-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    product_type VARCHAR(50) NOT NULL CHECK (product_type IN ('bateria', 'accesorio')),
    brand VARCHAR(100),
    model VARCHAR(100),
    voltage DECIMAL(10,2),
    amperage DECIMAL(10,2),
    battery_type VARCHAR(50) CHECK (battery_type IN ('seca', 'liquida')),
    polarity VARCHAR(50) CHECK (polarity IN ('izquierda', 'derecha')),
    acid_liters DECIMAL(10,3),
    vehicle_type VARCHAR(50) CHECK (vehicle_type IN ('auto', 'moto', 'otro')),
    min_sale_price DECIMAL(12,2) NOT NULL DEFAULT 0,
    effective_date TIMESTAMP WITH TIME ZONE,
    previous_price DECIMAL(12,2),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT
);

-- Create stock table
CREATE TABLE IF NOT EXISTS stock (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    product_type VARCHAR(50) NOT NULL,
    location_type VARCHAR(50) NOT NULL CHECK (location_type IN ('branch', 'warehouse')),
    location_id UUID NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0,
    min_stock_alert INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(product_id, location_type, location_id)
);

-- Create purchase_batches table
CREATE TABLE IF NOT EXISTS purchase_batches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    supplier_id UUID REFERENCES suppliers(id) ON DELETE SET NULL,
    invoice_number VARCHAR(100),
    purchase_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    notes TEXT,
    total_cost DECIMAL(12,2) NOT NULL DEFAULT 0,
    received BOOLEAN DEFAULT false,
    received_at TIMESTAMP WITH TIME ZONE,
    received_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create purchase_batch_details table
CREATE TABLE IF NOT EXISTS purchase_batch_details (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    purchase_batch_id UUID NOT NULL REFERENCES purchase_batches(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_cost DECIMAL(12,2) NOT NULL CHECK (unit_cost > 0),
    UNIQUE(purchase_batch_id, product_id)
);

-- Create transfers table
CREATE TABLE IF NOT EXISTS transfers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    origin_type VARCHAR(50) NOT NULL CHECK (origin_type IN ('branch', 'warehouse')),
    origin_id UUID NOT NULL,
    destination_type VARCHAR(50) NOT NULL CHECK (destination_type IN ('branch', 'warehouse')),
    destination_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pendiente' CHECK (status IN ('pendiente', 'aprobada', 'rechazada', 'enviada', 'recibida', 'cancelada')),
    requested_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    rejected_by UUID REFERENCES users(id) ON DELETE SET NULL,
    sent_by UUID REFERENCES users(id) ON DELETE SET NULL,
    received_by UUID REFERENCES users(id) ON DELETE SET NULL,
    notes TEXT,
    rejection_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create transfer_details table
CREATE TABLE IF NOT EXISTS transfer_details (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transfer_id UUID NOT NULL REFERENCES transfers(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    UNIQUE(transfer_id, product_id)
);

-- Create indexes for products
CREATE INDEX IF NOT EXISTS idx_products_type ON products(product_type);
CREATE INDEX IF NOT EXISTS idx_products_brand ON products(brand);
CREATE INDEX IF NOT EXISTS idx_products_active ON products(active);
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
CREATE INDEX IF NOT EXISTS idx_products_vehicle_type ON products(vehicle_type);

-- Create indexes for stock
CREATE INDEX IF NOT EXISTS idx_stock_product_id ON stock(product_id);
CREATE INDEX IF NOT EXISTS idx_stock_location ON stock(location_type, location_id);
CREATE INDEX IF NOT EXISTS idx_stock_quantity ON stock(quantity);

-- Create indexes for purchase_batches
CREATE INDEX IF NOT EXISTS idx_purchase_batches_supplier ON purchase_batches(supplier_id);
CREATE INDEX IF NOT EXISTS idx_purchase_batches_date ON purchase_batches(purchase_date);
CREATE INDEX IF NOT EXISTS idx_purchase_batches_received ON purchase_batches(received);

-- Create indexes for transfers
CREATE INDEX IF NOT EXISTS idx_transfers_origin ON transfers(origin_type, origin_id);
CREATE INDEX IF NOT EXISTS idx_transfers_destination ON transfers(destination_type, destination_id);
CREATE INDEX IF NOT EXISTS idx_transfers_status ON transfers(status);
CREATE INDEX IF NOT EXISTS idx_transfers_requested_by ON transfers(requested_by);

-- Add foreign key constraints for stock location (polymorphic reference)
-- Note: Since stock can reference either branches or warehouses, we use a check constraint
-- and the application layer ensures referential integrity
ALTER TABLE stock ADD CONSTRAINT chk_stock_location_type 
    CHECK (location_type IN ('branch', 'warehouse'));

-- Add foreign key constraints for transfer locations (polymorphic reference)
ALTER TABLE transfers ADD CONSTRAINT chk_transfer_origin_type 
    CHECK (origin_type IN ('branch', 'warehouse'));
ALTER TABLE transfers ADD CONSTRAINT chk_transfer_destination_type 
    CHECK (destination_type IN ('branch', 'warehouse'));

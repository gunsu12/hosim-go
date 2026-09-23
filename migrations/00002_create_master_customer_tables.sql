-- +goose Up
-- ========================================================
-- Modul: Master Customer / Penjamin (Customer & Customer Type)
-- ========================================================

CREATE TABLE IF NOT EXISTS customer_types (
    id VARCHAR(36) PRIMARY KEY,
    code VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    deleted_by VARCHAR(50)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_customer_types_code ON customer_types(code);
CREATE INDEX IF NOT EXISTS idx_customer_types_deleted_at ON customer_types(deleted_at);

CREATE TABLE IF NOT EXISTS customers (
    id VARCHAR(36) PRIMARY KEY,
    code VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    address VARCHAR(255),
    phone VARCHAR(255),
    email VARCHAR(255),
    website VARCHAR(255),
    contact_person VARCHAR(255),
    require_card BOOLEAN DEFAULT TRUE,
    description VARCHAR(255),
    customer_type_id VARCHAR(36) REFERENCES customer_types(id) ON DELETE SET NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_immutable BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    deleted_by VARCHAR(50)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_customers_code ON customers(code);
CREATE INDEX IF NOT EXISTS idx_customers_customer_type_id ON customers(customer_type_id);
CREATE INDEX IF NOT EXISTS idx_customers_is_active ON customers(is_active);
CREATE INDEX IF NOT EXISTS idx_customers_deleted_at ON customers(deleted_at);

-- +goose Down
DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS customer_types;

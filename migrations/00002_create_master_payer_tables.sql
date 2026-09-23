-- +goose Up
-- ========================================================
-- Modul: Master Penjamin (Payer & Payer Type)
-- ========================================================

CREATE TABLE IF NOT EXISTS payer_types (
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

CREATE UNIQUE INDEX IF NOT EXISTS idx_payer_types_code ON payer_types(code);
CREATE INDEX IF NOT EXISTS idx_payer_types_deleted_at ON payer_types(deleted_at);

CREATE TABLE IF NOT EXISTS payers (
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
    payer_type_id VARCHAR(36) REFERENCES payer_types(id) ON DELETE SET NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_immutable BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    deleted_by VARCHAR(50)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_payers_code ON payers(code);
CREATE INDEX IF NOT EXISTS idx_payers_payer_type_id ON payers(payer_type_id);
CREATE INDEX IF NOT EXISTS idx_payers_is_active ON payers(is_active);
CREATE INDEX IF NOT EXISTS idx_payers_deleted_at ON payers(deleted_at);

-- +goose Down
DROP TABLE IF EXISTS payers;
DROP TABLE IF EXISTS payer_types;

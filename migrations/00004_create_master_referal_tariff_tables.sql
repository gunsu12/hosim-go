-- +goose Up
-- ========================================================
-- Modul: Master Rujukan & Kelas Tarif
-- ========================================================

CREATE TABLE IF NOT EXISTS referals (
    id VARCHAR(36) PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    type VARCHAR(20) NOT NULL,
    address VARCHAR(255),
    phone VARCHAR(20),
    email VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    deleted_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM'
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_referals_code ON referals(code);
CREATE INDEX IF NOT EXISTS idx_referals_deleted_at ON referals(deleted_at);

CREATE TABLE IF NOT EXISTS tariff_classes (
    id VARCHAR(36) PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    description VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    deleted_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM'
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_tariff_classes_code ON tariff_classes(code);
CREATE INDEX IF NOT EXISTS idx_tariff_classes_is_active ON tariff_classes(is_active);
CREATE INDEX IF NOT EXISTS idx_tariff_classes_deleted_at ON tariff_classes(deleted_at);

-- +goose Down
DROP TABLE IF EXISTS tariff_classes;
DROP TABLE IF EXISTS referals;

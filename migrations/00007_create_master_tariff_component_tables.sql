-- +goose Up
-- ========================================================
-- Modul: Master Komponen Tarif (Tariff Component)
-- ========================================================

CREATE TABLE IF NOT EXISTS tariff_components (
    id VARCHAR(36) PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    description VARCHAR(255),
    is_hospital_revenue BOOLEAN NOT NULL DEFAULT FALSE,
    is_operator_revenue BOOLEAN NOT NULL DEFAULT FALSE,
    is_paramedic_revenue BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    deleted_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM'
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_tariff_components_code ON tariff_components(code);
CREATE INDEX IF NOT EXISTS idx_tariff_components_deleted_at ON tariff_components(deleted_at);

-- +goose Down
DROP TABLE IF EXISTS tariff_components;

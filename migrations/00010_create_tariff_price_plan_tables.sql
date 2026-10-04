-- +goose Up
-- ========================================================
-- Modul: Master Buku Tarif & Aturan Harga (Tariff Price Plan)
-- ========================================================

-- Perbarui tabel tariff_components jika kolom baru belum ada
ALTER TABLE tariff_components ADD COLUMN IF NOT EXISTS component_type VARCHAR(30) DEFAULT 'LAINNYA';
ALTER TABLE tariff_components ADD COLUMN IF NOT EXISTS default_coa_code VARCHAR(50);
ALTER TABLE tariff_components ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;

CREATE INDEX IF NOT EXISTS idx_tariff_components_is_active ON tariff_components(is_active);

-- Header Buku Tarif (tariff_price_plans)
CREATE TABLE IF NOT EXISTS tariff_price_plans (
    id VARCHAR(36) PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    effective_from DATE NOT NULL,
    effective_to DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    customer_id VARCHAR(36) REFERENCES customers(id) ON DELETE SET NULL,
    default_cito_percent NUMERIC(5,2) NOT NULL DEFAULT 25.00,
    approved_at TIMESTAMPTZ,
    approved_by VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    deleted_by VARCHAR(50)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_tariff_price_plans_code ON tariff_price_plans(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tariff_price_plans_status ON tariff_price_plans(status);
CREATE INDEX IF NOT EXISTS idx_tariff_price_plans_customer_id ON tariff_price_plans(customer_id);
CREATE INDEX IF NOT EXISTS idx_tariff_price_plans_is_default ON tariff_price_plans(is_default);
CREATE INDEX IF NOT EXISTS idx_tariff_price_plans_dates ON tariff_price_plans(effective_from, effective_to);
CREATE INDEX IF NOT EXISTS idx_tariff_price_plans_deleted_at ON tariff_price_plans(deleted_at);

-- Header Tarif Tindakan per Kelas (tariff_price_plan_items)
CREATE TABLE IF NOT EXISTS tariff_price_plan_items (
    id VARCHAR(36) PRIMARY KEY,
    price_plan_id VARCHAR(36) NOT NULL REFERENCES tariff_price_plans(id) ON DELETE CASCADE,
    item_id VARCHAR(36) NOT NULL REFERENCES items(id) ON DELETE RESTRICT,
    tariff_class_id VARCHAR(36) NOT NULL REFERENCES tariff_classes(id) ON DELETE RESTRICT,
    total_base_price NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    total_cito_price NUMERIC(15,2),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    deleted_by VARCHAR(50)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_tariff_price_plan_items_unique ON tariff_price_plan_items(price_plan_id, item_id, tariff_class_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tariff_price_plan_items_plan_id ON tariff_price_plan_items(price_plan_id);
CREATE INDEX IF NOT EXISTS idx_tariff_price_plan_items_item_id ON tariff_price_plan_items(item_id);
CREATE INDEX IF NOT EXISTS idx_tariff_price_plan_items_class_id ON tariff_price_plan_items(tariff_class_id);
CREATE INDEX IF NOT EXISTS idx_tariff_price_plan_items_deleted_at ON tariff_price_plan_items(deleted_at);

-- Rincian Pemecahan Komponen Biaya (tariff_price_plan_item_components)
CREATE TABLE IF NOT EXISTS tariff_price_plan_item_components (
    id VARCHAR(36) PRIMARY KEY,
    plan_item_id VARCHAR(36) NOT NULL REFERENCES tariff_price_plan_items(id) ON DELETE CASCADE,
    component_id VARCHAR(36) NOT NULL REFERENCES tariff_components(id) ON DELETE RESTRICT,
    base_amount NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    cito_amount NUMERIC(15,2),
    coa_code VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM'
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_tariff_price_plan_item_components_unique ON tariff_price_plan_item_components(plan_item_id, component_id);
CREATE INDEX IF NOT EXISTS idx_tariff_price_plan_item_components_component_id ON tariff_price_plan_item_components(component_id);

-- +goose Down
DROP TABLE IF EXISTS tariff_price_plan_item_components;
DROP TABLE IF EXISTS tariff_price_plan_items;
DROP TABLE IF EXISTS tariff_price_plans;
ALTER TABLE tariff_components DROP COLUMN IF EXISTS is_active;
ALTER TABLE tariff_components DROP COLUMN IF EXISTS default_coa_code;
ALTER TABLE tariff_components DROP COLUMN IF EXISTS component_type;

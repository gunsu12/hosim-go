-- +goose Up
-- ========================================================
-- Modul: Master Item Katalog (Class Table Inheritance)
-- ========================================================

CREATE TABLE IF NOT EXISTS item_categories (
    id VARCHAR(36) PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    item_type VARCHAR(20) NOT NULL,
    income_coa_code VARCHAR(50),
    discount_coa_code VARCHAR(50),
    sales_tax_coa_code VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    deleted_by VARCHAR(50)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_item_categories_code ON item_categories(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_item_categories_item_type ON item_categories(item_type);
CREATE INDEX IF NOT EXISTS idx_item_categories_deleted_at ON item_categories(deleted_at);

CREATE TABLE IF NOT EXISTS item_product_lines (
    id VARCHAR(36) PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    inventory_coa_code VARCHAR(50),
    cogs_coa_code VARCHAR(50),
    purchase_discount_coa_code VARCHAR(50),
    purchase_tax_coa_code VARCHAR(50),
    asset_coa_code VARCHAR(50),
    asset_accumulation_coa_code VARCHAR(50),
    asset_depreciation_expense_coa_code VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    deleted_by VARCHAR(50)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_item_product_lines_code ON item_product_lines(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_item_product_lines_deleted_at ON item_product_lines(deleted_at);

CREATE TABLE IF NOT EXISTS items (
    id VARCHAR(36) PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    generic_name VARCHAR(200),
    item_type VARCHAR(20) NOT NULL,
    category_id VARCHAR(36) NOT NULL REFERENCES item_categories(id) ON DELETE RESTRICT,
    product_line_id VARCHAR(36) REFERENCES item_product_lines(id) ON DELETE RESTRICT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    deleted_by VARCHAR(50)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_items_code ON items(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_items_name ON items(name);
CREATE INDEX IF NOT EXISTS idx_items_item_type ON items(item_type);
CREATE INDEX IF NOT EXISTS idx_items_category_id ON items(category_id);
CREATE INDEX IF NOT EXISTS idx_items_product_line_id ON items(product_line_id);
CREATE INDEX IF NOT EXISTS idx_items_is_active ON items(is_active);
CREATE INDEX IF NOT EXISTS idx_items_deleted_at ON items(deleted_at);

CREATE TABLE IF NOT EXISTS item_medications (
    item_id VARCHAR(36) PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,
    kfa_code VARCHAR(50),
    bpom_nie VARCHAR(50),
    dosage_form VARCHAR(50),
    strength_amount VARCHAR(50),
    strength_unit VARCHAR(20),
    default_route VARCHAR(50),
    medication_type VARCHAR(30),
    is_high_alert BOOLEAN NOT NULL DEFAULT FALSE,
    is_lasa BOOLEAN NOT NULL DEFAULT FALSE,
    is_fornas BOOLEAN NOT NULL DEFAULT FALSE,
    is_antibiotic BOOLEAN NOT NULL DEFAULT FALSE,
    storage_temperature VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM'
);

CREATE INDEX IF NOT EXISTS idx_item_medications_kfa_code ON item_medications(kfa_code);
CREATE INDEX IF NOT EXISTS idx_item_medications_bpom_nie ON item_medications(bpom_nie);

CREATE TABLE IF NOT EXISTS item_generals (
    item_id VARCHAR(36) PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,
    general_type VARCHAR(30) NOT NULL,
    is_sterile BOOLEAN NOT NULL DEFAULT FALSE,
    is_disposable BOOLEAN NOT NULL DEFAULT TRUE,
    is_cssd_item BOOLEAN NOT NULL DEFAULT FALSE,
    sterilization_method VARCHAR(30),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM'
);

CREATE INDEX IF NOT EXISTS idx_item_generals_general_type ON item_generals(general_type);
CREATE INDEX IF NOT EXISTS idx_item_generals_is_cssd_item ON item_generals(is_cssd_item);

CREATE TABLE IF NOT EXISTS item_assets (
    item_id VARCHAR(36) PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,
    brand VARCHAR(100),
    model_name VARCHAR(100),
    is_medical_equipment BOOLEAN NOT NULL DEFAULT FALSE,
    expected_life_years INT,
    depreciation_method VARCHAR(30),
    maintenance_interval_days INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM'
);

CREATE TABLE IF NOT EXISTS item_tariffs (
    item_id VARCHAR(36) PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,
    name_alias VARCHAR(200),
    charge_type VARCHAR(30) NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM'
);

CREATE INDEX IF NOT EXISTS idx_item_tariffs_charge_type ON item_tariffs(charge_type);

CREATE TABLE IF NOT EXISTS item_units (
    id VARCHAR(36) PRIMARY KEY,
    item_id VARCHAR(36) NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    unit_name VARCHAR(50) NOT NULL,
    conversion_factor NUMERIC(12,4) NOT NULL DEFAULT 1.0,
    is_base_unit BOOLEAN NOT NULL DEFAULT FALSE,
    is_purchase_unit BOOLEAN NOT NULL DEFAULT FALSE,
    is_dispense_unit BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM'
);

CREATE INDEX IF NOT EXISTS idx_item_units_item_id ON item_units(item_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_item_units_item_unit_name ON item_units(item_id, unit_name);

-- +goose Down
DROP TABLE IF EXISTS item_units;
DROP TABLE IF EXISTS item_tariffs;
DROP TABLE IF EXISTS item_assets;
DROP TABLE IF EXISTS item_generals;
DROP TABLE IF EXISTS item_medications;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS item_product_lines;
DROP TABLE IF EXISTS item_categories;

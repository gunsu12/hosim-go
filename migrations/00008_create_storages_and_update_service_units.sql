-- +goose Up
-- ========================================================
-- Modul: Master Storage (Gudang/Depo) & Relasi Service Unit
-- ========================================================

CREATE TABLE IF NOT EXISTS storages (
    id VARCHAR(36) PRIMARY KEY,
    departement_id VARCHAR(36) REFERENCES departements(id) ON DELETE SET NULL,
    name VARCHAR(150) NOT NULL,
    type VARCHAR(20) NOT NULL,
    capacity NUMERIC(15,2) NOT NULL DEFAULT 0,
    storage_parent_id VARCHAR(36) REFERENCES storages(id) ON DELETE SET NULL,
    location VARCHAR(150) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    deleted_by VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS idx_storages_departement_id ON storages(departement_id);
CREATE INDEX IF NOT EXISTS idx_storages_storage_parent_id ON storages(storage_parent_id);
CREATE INDEX IF NOT EXISTS idx_storages_deleted_at ON storages(deleted_at);

-- Update tabel service_units menambahkan kolom type, storage_id, dan default_class_id
ALTER TABLE service_units ADD COLUMN IF NOT EXISTS type VARCHAR(20) NOT NULL DEFAULT 'medical';
ALTER TABLE service_units ADD COLUMN IF NOT EXISTS storage_id VARCHAR(36) REFERENCES storages(id) ON DELETE SET NULL;
ALTER TABLE service_units ADD COLUMN IF NOT EXISTS default_class_id VARCHAR(36) REFERENCES tariff_classes(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_service_units_type ON service_units(type);
CREATE INDEX IF NOT EXISTS idx_service_units_storage_id ON service_units(storage_id);
CREATE INDEX IF NOT EXISTS idx_service_units_default_class_id ON service_units(default_class_id);

-- +goose Down
ALTER TABLE service_units DROP COLUMN IF EXISTS default_class_id;
ALTER TABLE service_units DROP COLUMN IF EXISTS storage_id;
ALTER TABLE service_units DROP COLUMN IF EXISTS type;
DROP TABLE IF EXISTS storages;

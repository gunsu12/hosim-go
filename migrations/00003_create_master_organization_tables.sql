-- +goose Up
-- ========================================================
-- Modul: Master Departemen / Instalasi, Unit Layanan & Ruangan
-- ========================================================

CREATE TABLE IF NOT EXISTS departements (
    id VARCHAR(36) PRIMARY KEY,
    code VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    address VARCHAR(255),
    phone VARCHAR(255),
    email VARCHAR(255),
    website VARCHAR(255),
    description VARCHAR(255),
    departement_type VARCHAR(20),
    ihs_organization_id VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    deleted_by VARCHAR(50)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_departements_code ON departements(code);
CREATE INDEX IF NOT EXISTS idx_departements_type ON departements(departement_type);
CREATE INDEX IF NOT EXISTS idx_departements_is_active ON departements(is_active);
CREATE INDEX IF NOT EXISTS idx_departements_deleted_at ON departements(deleted_at);

CREATE TABLE IF NOT EXISTS service_units (
    id VARCHAR(36) PRIMARY KEY,
    code VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    departement_id VARCHAR(36) REFERENCES departements(id) ON DELETE SET NULL,
    address VARCHAR(255),
    phone VARCHAR(255),
    email VARCHAR(255),
    website VARCHAR(255),
    description VARCHAR(255),
    ihs_location_id VARCHAR(255),
    is_registration_target BOOLEAN NOT NULL DEFAULT TRUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    deleted_by VARCHAR(50)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_service_units_code ON service_units(code);
CREATE INDEX IF NOT EXISTS idx_service_units_departement_id ON service_units(departement_id);
CREATE INDEX IF NOT EXISTS idx_service_units_is_reg_target ON service_units(is_registration_target);
CREATE INDEX IF NOT EXISTS idx_service_units_is_active ON service_units(is_active);
CREATE INDEX IF NOT EXISTS idx_service_units_deleted_at ON service_units(deleted_at);

CREATE TABLE IF NOT EXISTS rooms (
    id VARCHAR(36) PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    capacity INT NOT NULL DEFAULT 1,
    location VARCHAR(150),
    room_type VARCHAR(50) NOT NULL,
    service_unit_id VARCHAR(36) NOT NULL REFERENCES service_units(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    deleted_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM'
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_rooms_code ON rooms(code);
CREATE INDEX IF NOT EXISTS idx_rooms_service_unit_id ON rooms(service_unit_id);
CREATE INDEX IF NOT EXISTS idx_rooms_deleted_at ON rooms(deleted_at);

-- +goose Down
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS service_units;
DROP TABLE IF EXISTS departements;

-- +goose Up
-- ========================================================
-- Modul: Master Tenaga Medis / Practitioner (SatuSehat)
-- ========================================================

CREATE TABLE IF NOT EXISTS professions (
    id VARCHAR(36) PRIMARY KEY,
    code VARCHAR(20) NOT NULL,
    name VARCHAR(150) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    deleted_by VARCHAR(50)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_professions_code ON professions(code);
CREATE INDEX IF NOT EXISTS idx_professions_deleted_at ON professions(deleted_at);

CREATE TABLE IF NOT EXISTS specialties (
    id VARCHAR(36) PRIMARY KEY,
    code VARCHAR(20) NOT NULL,
    name VARCHAR(150) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    deleted_by VARCHAR(50)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_specialties_code ON specialties(code);
CREATE INDEX IF NOT EXISTS idx_specialties_deleted_at ON specialties(deleted_at);

CREATE TABLE IF NOT EXISTS practitioners (
    id VARCHAR(36) PRIMARY KEY,
    nik VARCHAR(16),
    nip VARCHAR(50),
    name VARCHAR(150) NOT NULL,
    gender VARCHAR(10) NOT NULL,
    sip VARCHAR(50),
    sip_expiry_date TIMESTAMPTZ,
    str VARCHAR(50),
    profession_id VARCHAR(36) REFERENCES professions(id) ON DELETE SET NULL,
    specialty_id VARCHAR(36) REFERENCES specialties(id) ON DELETE SET NULL,
    phone VARCHAR(20),
    email VARCHAR(100),
    ihs_practitioner_id VARCHAR(50),
    ihs_practitioner_name VARCHAR(150),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    deleted_by VARCHAR(50)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_practitioners_nik ON practitioners(nik)
WHERE nik IS NOT NULL AND nik != '' AND deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_practitioners_nip ON practitioners(nip)
WHERE nip IS NOT NULL AND nip != '' AND deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_practitioners_ihs ON practitioners(ihs_practitioner_id)
WHERE ihs_practitioner_id IS NOT NULL AND ihs_practitioner_id != '' AND deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_practitioners_profession_id ON practitioners(profession_id);
CREATE INDEX IF NOT EXISTS idx_practitioners_specialty_id ON practitioners(specialty_id);
CREATE INDEX IF NOT EXISTS idx_practitioners_is_active ON practitioners(is_active);
CREATE INDEX IF NOT EXISTS idx_practitioners_deleted_at ON practitioners(deleted_at);

-- +goose Down
DROP TABLE IF EXISTS practitioners;
DROP TABLE IF EXISTS specialties;
DROP TABLE IF EXISTS professions;

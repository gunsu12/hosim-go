-- +goose Up
-- ========================================================
-- Modul: Autentikasi & Pemetaan Tenaga Medis (Practitioner Link)
-- Menambahkan kolom practitioner_id ke tabel users (1 User <-> 1 Practitioner)
-- ========================================================

ALTER TABLE users ADD COLUMN IF NOT EXISTS practitioner_id VARCHAR(36) REFERENCES practitioners(id) ON DELETE RESTRICT;

CREATE UNIQUE INDEX IF NOT EXISTS uq_users_practitioner_id ON users(practitioner_id)
WHERE practitioner_id IS NOT NULL AND deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS uq_users_practitioner_id;
ALTER TABLE users DROP COLUMN IF EXISTS practitioner_id;

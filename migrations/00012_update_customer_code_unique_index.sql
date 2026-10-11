-- +goose Up
-- ========================================================
-- Perbarui unique index code customer menjadi partial index
-- Memungkinkan penggunaan ulang kode penjamin jika data lama di-soft delete
-- ========================================================

DROP INDEX IF EXISTS idx_customers_code;
CREATE UNIQUE INDEX IF NOT EXISTS idx_customers_code ON customers(code) WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_customers_code;
CREATE UNIQUE INDEX IF NOT EXISTS idx_customers_code ON customers(code);

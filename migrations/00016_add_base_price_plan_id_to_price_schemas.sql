-- +goose Up
-- ========================================================
-- Modul: Finance - Hubungkan Price Schema ke Buku Tarif Acuan (Price Plan)
-- ========================================================

ALTER TABLE price_schemas 
ADD COLUMN IF NOT EXISTS base_price_plan_id VARCHAR(36) REFERENCES tariff_price_plans(id) ON DELETE RESTRICT;

CREATE INDEX IF NOT EXISTS idx_price_schemas_base_price_plan_id ON price_schemas(base_price_plan_id);

-- +goose Down
DROP INDEX IF EXISTS idx_price_schemas_base_price_plan_id;
ALTER TABLE price_schemas DROP COLUMN IF EXISTS base_price_plan_id;

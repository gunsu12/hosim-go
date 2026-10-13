-- +goose Up
-- ========================================================
-- Modul: Finance - Skema Penyesuaian Harga PKS (Price Schema)
-- ========================================================

-- 1. Header Skema Harga (price_schemas)
CREATE TABLE IF NOT EXISTS price_schemas (
    id          VARCHAR(36)  PRIMARY KEY,
    code        VARCHAR(50)  NOT NULL,
    name        VARCHAR(150) NOT NULL,
    description TEXT,
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by  VARCHAR(50)  NOT NULL DEFAULT 'SYSTEM',
    updated_by  VARCHAR(50)  NOT NULL DEFAULT 'SYSTEM'
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_price_schemas_code ON price_schemas(code);
CREATE INDEX IF NOT EXISTS idx_price_schemas_is_active ON price_schemas(is_active);

-- 2. Aturan Bertingkat & Penyesuaian Tarif (price_schema_rules)
CREATE TABLE IF NOT EXISTS price_schema_rules (
    id                          VARCHAR(36)   PRIMARY KEY,
    price_schema_id             VARCHAR(36)   NOT NULL REFERENCES price_schemas(id) ON DELETE CASCADE,
    departement_id              VARCHAR(36)   REFERENCES departements(id) ON DELETE RESTRICT,
    service_unit_id             VARCHAR(36)   REFERENCES service_units(id) ON DELETE RESTRICT,
    item_id                     VARCHAR(36)   REFERENCES items(id) ON DELETE RESTRICT,

    -- Tarif tindakan (fase 1)
    tariff_adjustment_type      VARCHAR(20)   NOT NULL DEFAULT 'NONE',
    tariff_adjustment_value     NUMERIC(15,4) NOT NULL DEFAULT 0.0000,

    -- Obat & BHP (fase 2, disiapkan di skema)
    medication_adjustment_type  VARCHAR(20)   NOT NULL DEFAULT 'NONE',
    medication_adjustment_value NUMERIC(15,4) NOT NULL DEFAULT 0.0000,

    -- Cost sharing (% tanggungan dari total harga bruto setelah pajak)
    pay_by_patient_percent      NUMERIC(5,2)  NOT NULL DEFAULT 0.00,
    pay_by_customer_percent     NUMERIC(5,2)  NOT NULL DEFAULT 100.00,

    -- Pajak PPN (% dari final_price tarif)
    tax_percent                 NUMERIC(5,2)  NOT NULL DEFAULT 0.00,

    created_at                  TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by                  VARCHAR(50)   NOT NULL DEFAULT 'SYSTEM',
    updated_by                  VARCHAR(50)   NOT NULL DEFAULT 'SYSTEM',

    CONSTRAINT chk_psr_scope_exclusive
        CHECK (NOT (departement_id IS NOT NULL AND service_unit_id IS NOT NULL)),
    CONSTRAINT chk_psr_tariff_type
        CHECK (tariff_adjustment_type IN ('NONE', 'PERCENT', 'FIXED_AMOUNT', 'FIXED_PRICE')),
    CONSTRAINT chk_psr_medication_type
        CHECK (medication_adjustment_type IN ('NONE', 'PERCENT', 'FIXED_AMOUNT', 'FIXED_PRICE')),
    CONSTRAINT chk_psr_tariff_value CHECK (
        (tariff_adjustment_type = 'NONE'        AND tariff_adjustment_value = 0) OR
        (tariff_adjustment_type = 'PERCENT'     AND tariff_adjustment_value >= -100) OR
        (tariff_adjustment_type = 'FIXED_AMOUNT') OR
        (tariff_adjustment_type = 'FIXED_PRICE' AND tariff_adjustment_value >= 0 AND item_id IS NOT NULL)
    ),
    CONSTRAINT chk_psr_cost_share CHECK (
        pay_by_patient_percent BETWEEN 0 AND 100 AND
        pay_by_customer_percent BETWEEN 0 AND 100 AND
        pay_by_patient_percent + pay_by_customer_percent = 100
    ),
    CONSTRAINT chk_psr_tax CHECK (tax_percent BETWEEN 0 AND 100)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_psr_scope_unique ON price_schema_rules (
    price_schema_id,
    COALESCE(departement_id, ''),
    COALESCE(service_unit_id, ''),
    COALESCE(item_id, '')
);

CREATE INDEX IF NOT EXISTS idx_psr_lookup_item ON price_schema_rules(price_schema_id, item_id);
CREATE INDEX IF NOT EXISTS idx_psr_departement ON price_schema_rules(departement_id);
CREATE INDEX IF NOT EXISTS idx_psr_service_unit ON price_schema_rules(service_unit_id);

-- +goose Down
DROP TABLE IF EXISTS price_schema_rules;
DROP TABLE IF EXISTS price_schemas;

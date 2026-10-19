-- +goose Up
-- ========================================================
-- Modul: Kalender Akuntansi & Periode Fiskal (Fiscal Period)
-- ========================================================

CREATE TABLE IF NOT EXISTS fiscal_years (
    id VARCHAR(36) PRIMARY KEY,
    year INTEGER NOT NULL UNIQUE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_closed BOOLEAN NOT NULL DEFAULT FALSE,
    closed_at TIMESTAMPTZ,
    closed_by VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    CONSTRAINT chk_fiscal_years_dates CHECK (start_date <= end_date)
);

CREATE TABLE IF NOT EXISTS fiscal_periods (
    id VARCHAR(36) PRIMARY KEY,
    fiscal_year_id VARCHAR(36) NOT NULL REFERENCES fiscal_years(id) ON DELETE CASCADE,
    period_number INTEGER NOT NULL,
    period_key VARCHAR(7) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    locked_at TIMESTAMPTZ,
    locked_by VARCHAR(50),
    closed_at TIMESTAMPTZ,
    closed_by VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    CONSTRAINT chk_fiscal_periods_status CHECK (status IN ('OPEN', 'LOCKED', 'CLOSED')),
    CONSTRAINT chk_fiscal_periods_number CHECK (period_number BETWEEN 1 AND 13),
    CONSTRAINT chk_fiscal_periods_dates CHECK (start_date <= end_date)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_fiscal_years_year ON fiscal_years(year);
CREATE UNIQUE INDEX IF NOT EXISTS idx_fiscal_periods_period_key ON fiscal_periods(period_key);
CREATE INDEX IF NOT EXISTS idx_fiscal_periods_year_id ON fiscal_periods(fiscal_year_id);
CREATE INDEX IF NOT EXISTS idx_fiscal_periods_dates ON fiscal_periods(start_date, end_date);
CREATE INDEX IF NOT EXISTS idx_fiscal_periods_status ON fiscal_periods(status);

-- +goose Down
DROP TABLE IF EXISTS fiscal_periods;
DROP TABLE IF EXISTS fiscal_years;

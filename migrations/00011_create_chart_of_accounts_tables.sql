-- +goose Up
-- ========================================================
-- Modul: Bagan Akun / Chart of Accounts (COA)
-- ========================================================

CREATE TABLE IF NOT EXISTS chart_of_accounts (
    id VARCHAR(36) PRIMARY KEY,
    parent_id VARCHAR(36) REFERENCES chart_of_accounts(id) ON DELETE RESTRICT,
    code VARCHAR(50) NOT NULL,
    parent_code VARCHAR(50),
    name VARCHAR(150) NOT NULL,
    description TEXT,
    type VARCHAR(30) NOT NULL,
    position VARCHAR(10) NOT NULL,
    account_level INTEGER NOT NULL DEFAULT 1,
    is_postable BOOLEAN NOT NULL DEFAULT TRUE,
    is_treasury_account BOOLEAN NOT NULL DEFAULT FALSE,
    bank_name VARCHAR(100),
    bank_account_number VARCHAR(50),
    currency VARCHAR(10) NOT NULL DEFAULT 'IDR',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    deleted_by VARCHAR(50),
    CONSTRAINT chk_coa_position CHECK (position IN ('DEBIT', 'CREDIT')),
    CONSTRAINT chk_coa_type CHECK (type IN ('ASSET', 'LIABILITY', 'EQUITY', 'REVENUE', 'EXPENSE')),
    CONSTRAINT chk_coa_level CHECK (account_level >= 1)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_chart_of_accounts_code ON chart_of_accounts(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_chart_of_accounts_parent_id ON chart_of_accounts(parent_id);
CREATE INDEX IF NOT EXISTS idx_chart_of_accounts_type ON chart_of_accounts(type);
CREATE INDEX IF NOT EXISTS idx_chart_of_accounts_position ON chart_of_accounts(position);
CREATE INDEX IF NOT EXISTS idx_chart_of_accounts_is_postable ON chart_of_accounts(is_postable);
CREATE INDEX IF NOT EXISTS idx_chart_of_accounts_is_treasury ON chart_of_accounts(is_treasury_account);
CREATE INDEX IF NOT EXISTS idx_chart_of_accounts_is_active ON chart_of_accounts(is_active);
CREATE INDEX IF NOT EXISTS idx_chart_of_accounts_deleted_at ON chart_of_accounts(deleted_at);

-- +goose Down
DROP TABLE IF EXISTS chart_of_accounts;

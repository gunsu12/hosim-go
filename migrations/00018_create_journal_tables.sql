-- +goose Up
-- ========================================================
-- Modul: Pencatatan Jurnal & Buku Besar (Journal Entries)
-- ========================================================

CREATE TABLE IF NOT EXISTS journal_entries (
    id VARCHAR(36) PRIMARY KEY,
    entry_number VARCHAR(50) NOT NULL,
    entry_date DATE NOT NULL,
    fiscal_period_id VARCHAR(36) NOT NULL REFERENCES fiscal_periods(id) ON DELETE RESTRICT,
    fiscal_period VARCHAR(7) NOT NULL,
    entry_type VARCHAR(30) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    description TEXT NOT NULL,
    reference_number VARCHAR(100),
    total_debit NUMERIC(18,2) NOT NULL DEFAULT 0,
    total_credit NUMERIC(18,2) NOT NULL DEFAULT 0,
    reversed_entry_id VARCHAR(36) REFERENCES journal_entries(id) ON DELETE SET NULL,
    reversal_reason TEXT,
    posted_at TIMESTAMPTZ,
    posted_by VARCHAR(50),
    reversed_at TIMESTAMPTZ,
    reversed_by VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    updated_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    deleted_by VARCHAR(50),
    CONSTRAINT chk_journal_entries_type CHECK (entry_type IN ('GENERAL', 'ADJUSTMENT', 'CLOSING', 'REVERSAL', 'SYSTEM_BILLING', 'SYSTEM_CASHIER', 'SYSTEM_INVENTORY')),
    CONSTRAINT chk_journal_entries_status CHECK (status IN ('DRAFT', 'POSTED', 'REVERSED')),
    CONSTRAINT chk_journal_entries_amounts CHECK (total_debit >= 0 AND total_credit >= 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_journal_entries_entry_number ON journal_entries(entry_number) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_journal_entries_date_status ON journal_entries(entry_date, status);
CREATE INDEX IF NOT EXISTS idx_journal_entries_fiscal_period_id ON journal_entries(fiscal_period_id);
CREATE INDEX IF NOT EXISTS idx_journal_entries_fiscal_period ON journal_entries(fiscal_period);
CREATE INDEX IF NOT EXISTS idx_journal_entries_entry_type ON journal_entries(entry_type);
CREATE INDEX IF NOT EXISTS idx_journal_entries_deleted_at ON journal_entries(deleted_at);

CREATE TABLE IF NOT EXISTS journal_entry_lines (
    id VARCHAR(36) PRIMARY KEY,
    journal_entry_id VARCHAR(36) NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
    line_number INTEGER NOT NULL,
    account_id VARCHAR(36) NOT NULL REFERENCES chart_of_accounts(id) ON DELETE RESTRICT,
    account_code VARCHAR(50) NOT NULL,
    account_name VARCHAR(150) NOT NULL,
    entry_position VARCHAR(10) NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    memo TEXT,
    department_id VARCHAR(36),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_journal_lines_position CHECK (entry_position IN ('DEBIT', 'CREDIT')),
    CONSTRAINT chk_journal_lines_amount CHECK (amount > 0)
);

CREATE INDEX IF NOT EXISTS idx_journal_entry_lines_entry_id ON journal_entry_lines(journal_entry_id);
CREATE INDEX IF NOT EXISTS idx_journal_entry_lines_account_id ON journal_entry_lines(account_id);
CREATE INDEX IF NOT EXISTS idx_journal_entry_lines_account_code ON journal_entry_lines(account_code);

-- +goose Down
DROP TABLE IF EXISTS journal_entry_lines;
DROP TABLE IF EXISTS journal_entries;

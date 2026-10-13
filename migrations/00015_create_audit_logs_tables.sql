-- +goose Up
-- ========================================================
-- Modul: Audit - Jejak Audit Append-Only (Audit Logs)
-- ========================================================

CREATE TABLE IF NOT EXISTS audit_logs (
    id             VARCHAR(36)  PRIMARY KEY,
    actor_id       VARCHAR(50)  NOT NULL,
    actor_username VARCHAR(100),
    actor_role     VARCHAR(50)  NOT NULL,
    action         VARCHAR(50)  NOT NULL,
    resource_type  VARCHAR(50)  NOT NULL,
    resource_id    VARCHAR(36)  NOT NULL,
    before_payload JSONB,
    after_payload  JSONB,
    correlation_id VARCHAR(100),
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_actor ON audit_logs(actor_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);

-- +goose Down
DROP TABLE IF EXISTS audit_logs;

CREATE TABLE IF NOT EXISTS security_audit_events (
    id TEXT PRIMARY KEY,
    category TEXT NOT NULL,
    action TEXT NOT NULL,
    severity TEXT NOT NULL DEFAULT 'info',
    actor TEXT NOT NULL DEFAULT 'system',
    resource_id TEXT NOT NULL DEFAULT '',
    summary TEXT NOT NULL,
    metadata_json TEXT NOT NULL DEFAULT '{}',
    ip_address TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_security_audit_created_at ON security_audit_events(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_security_audit_category ON security_audit_events(category);
CREATE INDEX IF NOT EXISTS idx_security_audit_severity ON security_audit_events(severity);

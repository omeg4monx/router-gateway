CREATE TABLE IF NOT EXISTS system_notifications (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    severity TEXT NOT NULL DEFAULT 'info',
    category TEXT NOT NULL DEFAULT 'system',
    source TEXT NOT NULL DEFAULT 'system',
    metadata_json TEXT NOT NULL DEFAULT '{}',
    read_at TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON system_notifications(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_read_at ON system_notifications(read_at);
CREATE INDEX IF NOT EXISTS idx_notifications_severity ON system_notifications(severity);

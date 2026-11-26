CREATE TABLE IF NOT EXISTS plugin_configs (
    key         TEXT PRIMARY KEY,
    enabled     INTEGER NOT NULL DEFAULT 0,
    models_json TEXT NOT NULL DEFAULT '[]',
    config_json TEXT NOT NULL DEFAULT '{}',
    updated_at  TEXT NOT NULL
);

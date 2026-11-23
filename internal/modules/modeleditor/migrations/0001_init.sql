CREATE TABLE IF NOT EXISTS model_overrides (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    model_pattern TEXT NOT NULL,
    target_model TEXT NOT NULL DEFAULT '',
    system_prompt_prefix TEXT NOT NULL DEFAULT '',
    system_prompt_suffix TEXT NOT NULL DEFAULT '',
    temperature REAL,
    max_tokens INTEGER,
    top_p REAL,
    extra_params_json TEXT NOT NULL DEFAULT '{}',
    strip_params_json TEXT NOT NULL DEFAULT '[]',
    priority INTEGER NOT NULL DEFAULT 100,
    enabled INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_model_overrides_priority ON model_overrides(priority ASC, created_at ASC);

CREATE TABLE IF NOT EXISTS model_bypass_rules (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    model_pattern TEXT NOT NULL,
    key_id TEXT NOT NULL DEFAULT '',
    protocol TEXT NOT NULL DEFAULT '',
    header_name TEXT NOT NULL DEFAULT '',
    header_value TEXT NOT NULL DEFAULT '',
    action TEXT NOT NULL,
    target_model TEXT NOT NULL DEFAULT '',
    priority INTEGER NOT NULL DEFAULT 100,
    enabled INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_model_bypass_rules_priority ON model_bypass_rules(priority ASC, created_at ASC);

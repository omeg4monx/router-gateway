CREATE TABLE IF NOT EXISTS prompt_guard_config (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    enabled INTEGER NOT NULL DEFAULT 0,
    mode TEXT NOT NULL DEFAULT 'audit',
    sensitivity TEXT NOT NULL DEFAULT 'medium',
    block_status_code INTEGER NOT NULL DEFAULT 400,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS prompt_guard_rules (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    pattern TEXT NOT NULL,
    match_type TEXT NOT NULL,
    action TEXT NOT NULL,
    severity TEXT NOT NULL,
    enabled INTEGER NOT NULL DEFAULT 1,
    builtin INTEGER NOT NULL DEFAULT 0,
    description TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS prompt_guard_events (
    id TEXT PRIMARY KEY,
    request_id TEXT NOT NULL,
    model TEXT NOT NULL,
    key_id TEXT NOT NULL,
    rule_id TEXT NOT NULL,
    rule_name TEXT NOT NULL,
    action TEXT NOT NULL,
    severity TEXT NOT NULL,
    snippet TEXT NOT NULL,
    created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_prompt_guard_events_created_at ON prompt_guard_events(created_at DESC);

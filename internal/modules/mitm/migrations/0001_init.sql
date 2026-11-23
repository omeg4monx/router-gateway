CREATE TABLE IF NOT EXISTS mitm_config (
    id                 INTEGER PRIMARY KEY CHECK (id = 1),
    listen_addr        TEXT    NOT NULL DEFAULT '127.0.0.1:8443',
    target_gateway_url TEXT    NOT NULL DEFAULT '',
    enabled_tools_json TEXT    NOT NULL DEFAULT '["antigravity","copilot","kiro","cursor"]',
    auto_start         INTEGER NOT NULL DEFAULT 0,
    updated_at         TEXT    NOT NULL
);

CREATE TABLE IF NOT EXISTS mitm_aliases (
    id           TEXT    PRIMARY KEY,
    tool         TEXT    NOT NULL,
    source_model TEXT    NOT NULL,
    target_model TEXT    NOT NULL,
    enabled      INTEGER NOT NULL DEFAULT 1,
    created_at   TEXT    NOT NULL,
    updated_at   TEXT    NOT NULL,
    UNIQUE(tool, source_model)
);

CREATE INDEX IF NOT EXISTS idx_mitm_aliases_tool ON mitm_aliases(tool);

CREATE TABLE IF NOT EXISTS mitm_events (
    id           TEXT    PRIMARY KEY,
    tool         TEXT    NOT NULL,
    host         TEXT    NOT NULL,
    method       TEXT    NOT NULL,
    path         TEXT    NOT NULL,
    source_model TEXT    NOT NULL DEFAULT '',
    target_model TEXT    NOT NULL DEFAULT '',
    status_code  INTEGER NOT NULL DEFAULT 0,
    latency_ms   INTEGER NOT NULL DEFAULT 0,
    created_at   TEXT    NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_mitm_events_created_at ON mitm_events(created_at DESC);

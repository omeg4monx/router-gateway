CREATE TABLE IF NOT EXISTS tokensaver_config (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    enabled INTEGER NOT NULL DEFAULT 0,
    strategy TEXT NOT NULL DEFAULT 'whitespace',
    max_history_messages INTEGER NOT NULL DEFAULT 40,
    max_tool_output_chars INTEGER NOT NULL DEFAULT 8000,
    headroom_enabled INTEGER NOT NULL DEFAULT 1,
    prefix_depth INTEGER NOT NULL DEFAULT 2,
    max_cache_entries INTEGER NOT NULL DEFAULT 1024,
    cache_ttl_seconds INTEGER NOT NULL DEFAULT 1800,
    inject_cache_control INTEGER NOT NULL DEFAULT 1,
    updated_at TEXT NOT NULL
);

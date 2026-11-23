CREATE TABLE IF NOT EXISTS telegram_config (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    enabled INTEGER NOT NULL DEFAULT 0,
    bot_token_enc TEXT NOT NULL DEFAULT '',
    bot_token_mask TEXT NOT NULL DEFAULT '',
    chat_id TEXT NOT NULL DEFAULT '',
    user_id TEXT NOT NULL DEFAULT '',
    webhook_url TEXT NOT NULL DEFAULT '',
    webhook_secret_enc TEXT NOT NULL DEFAULT '',
    notify_on_job_failure INTEGER NOT NULL DEFAULT 1,
    notify_on_auth_change INTEGER NOT NULL DEFAULT 1,
    notify_on_emit INTEGER NOT NULL DEFAULT 0,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS telegram_events (
    id TEXT PRIMARY KEY,
    direction TEXT NOT NULL,
    kind TEXT NOT NULL,
    command TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL,
    summary TEXT NOT NULL DEFAULT '',
    error TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_telegram_events_created_at
    ON telegram_events(created_at DESC);

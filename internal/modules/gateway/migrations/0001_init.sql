CREATE TABLE IF NOT EXISTS response_chains (
    response_id TEXT PRIMARY KEY,
    key_id TEXT NOT NULL,
    model TEXT NOT NULL,
    messages_json TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_response_chains_expires
    ON response_chains(expires_at);

CREATE TABLE IF NOT EXISTS vouchers (
    id TEXT PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    quota_usd REAL NOT NULL DEFAULT 0,
    rate_limit_rpm INTEGER NOT NULL DEFAULT 0,
    allowed_models_json TEXT NOT NULL DEFAULT '[]',
    team_id TEXT NOT NULL DEFAULT '',
    max_redemptions INTEGER NOT NULL DEFAULT 1,
    redeemed_count INTEGER NOT NULL DEFAULT 0,
    enabled INTEGER NOT NULL DEFAULT 1,
    expires_at TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS voucher_redemptions (
    id TEXT PRIMARY KEY,
    voucher_id TEXT NOT NULL,
    code TEXT NOT NULL,
    claimer_id TEXT NOT NULL,
    key_id TEXT NOT NULL DEFAULT '',
    ip_address TEXT NOT NULL DEFAULT '',
    redeemed_at TEXT NOT NULL,
    FOREIGN KEY (voucher_id) REFERENCES vouchers(id) ON DELETE CASCADE,
    UNIQUE(voucher_id, claimer_id)
);

CREATE INDEX IF NOT EXISTS idx_voucher_redemptions_voucher_id ON voucher_redemptions(voucher_id, redeemed_at DESC);

ALTER TABLE tokensaver_config ADD COLUMN caveman_enabled INTEGER NOT NULL DEFAULT 0;
ALTER TABLE tokensaver_config ADD COLUMN caveman_level TEXT NOT NULL DEFAULT 'full';
ALTER TABLE tokensaver_config ADD COLUMN ponytail_enabled INTEGER NOT NULL DEFAULT 0;
ALTER TABLE tokensaver_config ADD COLUMN ponytail_level TEXT NOT NULL DEFAULT 'full';
ALTER TABLE tokensaver_config ADD COLUMN trim_history_enabled INTEGER NOT NULL DEFAULT 0;

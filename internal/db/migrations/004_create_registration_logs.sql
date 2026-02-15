-- 004_create_registration_logs.sql

CREATE TABLE IF NOT EXISTS registration_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    profile_id INTEGER NOT NULL,
    triggered_by TEXT NOT NULL,
    started_at DATETIME NOT NULL,
    finished_at DATETIME,
    success INTEGER NOT NULL,
    error_code TEXT,
    error_message TEXT,
    logs TEXT,
    created_at DATETIME NOT NULL,
    CONSTRAINT fk_logs_profile FOREIGN KEY (profile_id)
        REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_logs_profile_created
    ON registration_logs (profile_id, created_at DESC);

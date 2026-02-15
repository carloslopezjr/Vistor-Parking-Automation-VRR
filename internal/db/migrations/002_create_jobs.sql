-- 002_create_jobs.sql

CREATE TABLE IF NOT EXISTS jobs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    job_type TEXT NOT NULL,
    profile_id INTEGER NOT NULL,
    run_at DATETIME NOT NULL,
    status TEXT NOT NULL,
    attempts INTEGER NOT NULL DEFAULT 0,
    last_error TEXT,
    created_at DATETIME NOT NULL,
    CONSTRAINT fk_jobs_profile FOREIGN KEY (profile_id)
        REFERENCES profiles(id) ON DELETE CASCADE,
    CONSTRAINT chk_job_status CHECK (status IN ('pending','running','completed','failed')),
    CONSTRAINT chk_job_type CHECK (job_type IN ('reminder_12h','reminder_20h','reminder_23h'))
);

CREATE INDEX IF NOT EXISTS idx_jobs_run_at_status ON jobs (status, run_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_jobs_unique_type_profile_run_at
    ON jobs (job_type, profile_id, run_at);

package scheduler

import (
	"context"
	"database/sql"
	"time"

	"vistor-parking-automation-vrr/internal/jobs"
	"vistor-parking-automation-vrr/internal/models"
)

// Logger is a minimal logging interface to avoid tying the scheduler to a
// specific logging implementation.
type Logger interface {
	Printf(format string, v ...any)
}

// Config controls scheduler behavior.
type Config struct {
	Interval          time.Duration // how often to poll for jobs (default 1m)
	BatchSize         int           // max jobs per run (default 10)
	StaleRunningAfter time.Duration // how long a job can stay 'running' before reset
	MaxAttempts       int           // max attempts before marking failed (default 5)
	BaseBackoff       time.Duration // base backoff duration (default 5m)
	MaxBackoff        time.Duration // max backoff duration (default 1h)
}

// Scheduler polls the jobs table and dispatches due jobs to the jobs service.
type Scheduler struct {
	db      *sql.DB
	jobsSvc jobs.Service
	log     Logger
	cfg     Config
}

// New constructs a new scheduler.
func New(db *sql.DB, jobsSvc jobs.Service, logger Logger, cfg Config) *Scheduler {
	if cfg.Interval <= 0 {
		cfg.Interval = time.Minute
	}
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 10
	}
	if cfg.StaleRunningAfter <= 0 {
		cfg.StaleRunningAfter = 10 * time.Minute
	}
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 5
	}
	if cfg.BaseBackoff <= 0 {
		cfg.BaseBackoff = 5 * time.Minute
	}
	if cfg.MaxBackoff <= 0 {
		cfg.MaxBackoff = time.Hour
	}

	return &Scheduler{db: db, jobsSvc: jobsSvc, log: logger, cfg: cfg}
}

// Run starts the scheduling loop and blocks until the context is cancelled or
// an unrecoverable error occurs.
func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.runOnce(ctx); err != nil {
				if s.log != nil {
					s.log.Printf("scheduler runOnce error: %v", err)
				}
			}
		}
	}
}

// RunOnce executes a single scheduling cycle; useful for tests.
func (s *Scheduler) RunOnce(ctx context.Context) error {
	return s.runOnce(ctx)
}

func (s *Scheduler) runOnce(ctx context.Context) error {
	now := time.Now().UTC()

	// Reset stale running jobs back to pending so they can be retried after a crash.
	if err := s.resetStaleRunning(ctx, now); err != nil {
		return err
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT id, job_type, profile_id, run_at, status, attempts, last_error, created_at
		 FROM jobs
		 WHERE status = 'pending' AND run_at <= ?
		 ORDER BY run_at
		 LIMIT ?`, now, s.cfg.BatchSize,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var job models.Job
		var lastError sql.NullString
		if err := rows.Scan(&job.ID, &job.JobType, &job.ProfileID, &job.RunAt, &job.Status, &job.Attempts, &lastError, &job.CreatedAt); err != nil {
			return err
		}
		if lastError.Valid {
			job.LastError = &lastError.String
		}

		claimed, err := s.claimJob(ctx, job.ID)
		if err != nil {
			return err
		}
		if !claimed {
			continue
		}

		// attempts was incremented in claimJob in the DB; reflect that locally.
		job.Attempts++

		if err := s.processJob(ctx, job, now); err != nil && s.log != nil {
			s.log.Printf("job %d processing error: %v", job.ID, err)
		}
	}
	return rows.Err()
}

func (s *Scheduler) resetStaleRunning(ctx context.Context, now time.Time) error {
	deadline := now.Add(-s.cfg.StaleRunningAfter)
	_, err := s.db.ExecContext(ctx,
		`UPDATE jobs
		 SET status = 'pending'
		 WHERE status = 'running' AND run_at <= ?`, deadline,
	)
	return err
}

func (s *Scheduler) claimJob(ctx context.Context, id int64) (bool, error) {
	res, err := s.db.ExecContext(ctx,
		`UPDATE jobs
		 SET status = 'running', attempts = attempts + 1
		 WHERE id = ? AND status = 'pending'`, id,
	)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected == 1, nil
}

func (s *Scheduler) processJob(ctx context.Context, job models.Job, now time.Time) error {
	// Dispatch to job service.
	err := s.jobsSvc.HandleJob(ctx, job)
	if err == nil {
		_, uerr := s.db.ExecContext(ctx,
			`UPDATE jobs SET status = 'completed', last_error = NULL WHERE id = ?`, job.ID,
		)
		return uerr
	}

	// Failure path with exponential backoff.
	attempts := job.Attempts
	if attempts >= s.cfg.MaxAttempts {
		_, uerr := s.db.ExecContext(ctx,
			`UPDATE jobs SET status = 'failed', last_error = ? WHERE id = ?`, err.Error(), job.ID,
		)
		if uerr != nil {
			return uerr
		}
		return err
	}

	backoff := s.cfg.BaseBackoff * (1 << (attempts - 1))
	if backoff > s.cfg.MaxBackoff {
		backoff = s.cfg.MaxBackoff
	}
	nextRun := now.Add(backoff)

	_, uerr := s.db.ExecContext(ctx,
		`UPDATE jobs
		 SET status = 'pending', run_at = ?, last_error = ?
		 WHERE id = ?`, nextRun, err.Error(), job.ID,
	)
	if uerr != nil {
		return uerr
	}

	return err
}

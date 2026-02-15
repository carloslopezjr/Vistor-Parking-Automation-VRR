package scheduler

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"vistor-parking-automation-vrr/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type dummyJobsService struct {
	handled []models.Job
}

func (d *dummyJobsService) ScheduleReminderJobs(ctx context.Context, profileID int64, baseTime time.Time) error {
	return nil
}

func (d *dummyJobsService) HandleJob(ctx context.Context, job models.Job) error {
	d.handled = append(d.handled, job)
	return nil
}

func newSchedulerTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		job_type TEXT NOT NULL,
		profile_id INTEGER NOT NULL,
		run_at DATETIME NOT NULL,
		status TEXT NOT NULL,
		attempts INTEGER NOT NULL DEFAULT 0,
		last_error TEXT,
		created_at DATETIME NOT NULL
	)`)
	if err != nil {
		t.Fatalf("create table: %v", err)
	}
	return db
}

func TestSchedulerDispatchesDueJob(t *testing.T) {
	db := newSchedulerTestDB(t)
	t.Cleanup(func() { db.Close() })

	now := time.Now().UTC()
	_, err := db.Exec(`INSERT INTO jobs (job_type, profile_id, run_at, status, attempts, created_at)
		VALUES (?, ?, ?, 'pending', 0, ?)`, models.JobTypeReminder12h, int64(1), now.Add(-time.Minute), now)
	if err != nil {
		t.Fatalf("insert job: %v", err)
	}

	dummy := &dummyJobsService{}
	s := New(db, dummy, nil, Config{
		Interval:    time.Hour,
		BatchSize:   5,
		MaxAttempts: 3,
		BaseBackoff: time.Minute,
		MaxBackoff:  time.Hour,
	})

	if err := s.RunOnce(context.Background()); err != nil {
		t.Fatalf("RunOnce error: %v", err)
	}

	if len(dummy.handled) != 1 {
		t.Fatalf("expected 1 handled job, got %d", len(dummy.handled))
	}
}

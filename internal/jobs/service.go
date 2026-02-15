package jobs

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"vistor-parking-automation-vrr/internal/models"
)

// ReminderHandler abstracts the work required to execute a reminder job.
// The concrete implementation will typically generate a token and send an
// email reminder.
type ReminderHandler interface {
	HandleReminder(ctx context.Context, job models.Job) error
}

// Service provides high-level job orchestration.
type Service interface {
	// ScheduleReminderJobs creates the 12h, 20h and 23h reminder jobs for a
	// profile, based on the given base time (typically the registration time).
	// The operation is idempotent: if jobs already exist for the same
	// profile/time/type combination, they are left as-is.
	ScheduleReminderJobs(ctx context.Context, profileID int64, baseTime time.Time) error

	// HandleJob dispatches a single job to the appropriate handler.
	HandleJob(ctx context.Context, job models.Job) error
}

// NewService constructs a DB-backed job service.
func NewService(db *sql.DB, reminderHandler ReminderHandler) Service {
	return &service{db: db, reminderHandler: reminderHandler}
}

type service struct {
	db              *sql.DB
	reminderHandler ReminderHandler
}

func (s *service) ScheduleReminderJobs(ctx context.Context, profileID int64, baseTime time.Time) error {
	if profileID <= 0 {
		return errors.New("invalid profile id")
	}

	base := baseTime.UTC()
	createdAt := time.Now().UTC()

	type jobSpec struct {
		jobType models.JobType
		runAt   time.Time
	}

	specs := []jobSpec{
		{jobType: models.JobTypeReminder12h, runAt: base.Add(12 * time.Hour)},
		{jobType: models.JobTypeReminder20h, runAt: base.Add(20 * time.Hour)},
		{jobType: models.JobTypeReminder23h, runAt: base.Add(23 * time.Hour)},
	}

	for _, spec := range specs {
		// INSERT OR IGNORE leverages the UNIQUE constraint on
		// (job_type, profile_id, run_at) to make this idempotent.
		_, err := s.db.ExecContext(ctx,
			`INSERT OR IGNORE INTO jobs (job_type, profile_id, run_at, status, attempts, created_at)
			 VALUES (?, ?, ?, 'pending', 0, ?)`,
			spec.jobType, profileID, spec.runAt, createdAt,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) HandleJob(ctx context.Context, job models.Job) error {
	switch job.JobType {
	case models.JobTypeReminder12h, models.JobTypeReminder20h, models.JobTypeReminder23h:
		if s.reminderHandler == nil {
			return errors.New("no reminder handler configured")
		}
		return s.reminderHandler.HandleReminder(ctx, job)
	default:
		return errors.New("unsupported job type")
	}
}

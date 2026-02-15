package models

import "time"

// Profile represents a parking profile stored in the database.
type Profile struct {
	ID                 int64
	ApartmentName      string
	LicensePlate       string
	VehicleMake        string
	VehicleModel       string
	ResidentName       string
	UnitNumber         string
	VisitorName        string
	ResidentEmail      string
	ProfileName        string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	LastRegistrationAt *time.Time
}

// JobType represents the type of scheduled job.
type JobType string

const (
	JobTypeReminder12h JobType = "reminder_12h"
	JobTypeReminder20h JobType = "reminder_20h"
	JobTypeReminder23h JobType = "reminder_23h"
)

// JobStatus represents the processing status of a job.
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

// Job represents a scheduled job in the system.
type Job struct {
	ID        int64
	JobType   JobType
	ProfileID int64
	RunAt     time.Time
	Status    JobStatus
	Attempts  int
	LastError *string
	CreatedAt time.Time
}

// ReminderToken represents a secure, single-use reminder token.
type ReminderToken struct {
	ID        int64
	ProfileID int64
	TokenHash []byte
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

// RegistrationTrigger identifies what initiated a registration attempt.
type RegistrationTrigger string

const (
	TriggerManualUI   RegistrationTrigger = "manual_ui"
	TriggerReminder12 RegistrationTrigger = "reminder_12h"
	TriggerReminder20 RegistrationTrigger = "reminder_20h"
	TriggerReminder23 RegistrationTrigger = "reminder_23h"
	TriggerTokenClick RegistrationTrigger = "token_click"
)

// RegistrationLog stores the outcome of an automation run.
type RegistrationLog struct {
	ID           int64
	ProfileID    int64
	TriggeredBy  RegistrationTrigger
	StartedAt    time.Time
	FinishedAt   *time.Time
	Success      bool
	ErrorCode    *string
	ErrorMessage *string
	Logs         *string
	CreatedAt    time.Time
}

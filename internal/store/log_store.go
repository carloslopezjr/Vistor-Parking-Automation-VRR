package store

import (
	"context"
	"database/sql"

	"vistor-parking-automation-vrr/internal/models"
)

// LogStore provides access to registration logs.
type LogStore struct {
	db *sql.DB
}

func NewLogStore(db *sql.DB) *LogStore {
	return &LogStore{db: db}
}

// LatestByProfile returns the latest log for each profile ID provided.
func (s *LogStore) LatestByProfile(ctx context.Context, profileIDs []int64) (map[int64]models.RegistrationLog, error) {
	res := make(map[int64]models.RegistrationLog)
	if len(profileIDs) == 0 {
		return res, nil
	}

	// Simple approach: query per profile for clarity; can be optimized later.
	for _, id := range profileIDs {
		var l models.RegistrationLog
		var finishedAt sql.NullTime
		var errorCode, errorMessage, logs sql.NullString
		err := s.db.QueryRowContext(ctx,
			`SELECT id, profile_id, triggered_by, started_at, finished_at, success,
			 error_code, error_message, logs, created_at
			 FROM registration_logs
			 WHERE profile_id = ?
			 ORDER BY created_at DESC
			 LIMIT 1`, id,
		).Scan(&l.ID, &l.ProfileID, &l.TriggeredBy, &l.StartedAt, &finishedAt, &l.Success,
			&errorCode, &errorMessage, &logs, &l.CreatedAt)
		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			return nil, err
		}
		if finishedAt.Valid {
			l.FinishedAt = &finishedAt.Time
		}
		if errorCode.Valid {
			l.ErrorCode = &errorCode.String
		}
		if errorMessage.Valid {
			l.ErrorMessage = &errorMessage.String
		}
		if logs.Valid {
			l.Logs = &logs.String
		}
		res[id] = l
	}
	return res, nil
}

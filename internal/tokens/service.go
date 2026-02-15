package tokens

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"strings"
	"time"
)

var (
	// ErrInvalidToken is returned when a token cannot be found, is expired, or already used.
	ErrInvalidToken = errors.New("invalid or expired token")
)

// Service defines operations for secure reminder tokens.
type Service interface {
	// Generate creates and stores a new single-use token for the given profile
	// and returns the raw token string for use in URLs. The token expires
	// automatically after the configured TTL (48h for reminders).
	Generate(ctx context.Context, profileID int64) (string, error)

	// ValidateAndConsume checks the token, ensures it is not expired and not
	// previously used, then marks it as used atomically and returns the
	// associated profile ID. If invalid, ErrInvalidToken is returned.
	ValidateAndConsume(ctx context.Context, rawToken string) (int64, error)
}

// Config contains settings for the token service.
type Config struct {
	// TTL is the lifetime of a token. For reminder links this should be 48h.
	TTL time.Duration
}

// NewService constructs a DB-backed token service.
func NewService(db *sql.DB, cfg Config) Service {
	if cfg.TTL <= 0 {
		cfg.TTL = 48 * time.Hour
	}
	return &service{db: db, cfg: cfg}
}

type service struct {
	db  *sql.DB
	cfg Config
}

func (s *service) Generate(ctx context.Context, profileID int64) (string, error) {
	// 32 bytes of entropy, URL-safe when base64url encoded.
	const tokenBytes = 32

	for {
		b := make([]byte, tokenBytes)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}

		token := base64.RawURLEncoding.EncodeToString(b)
		hash := sha256.Sum256([]byte(token))

		expiresAt := time.Now().Add(s.cfg.TTL).UTC()
		createdAt := time.Now().UTC()

		_, err := s.db.ExecContext(ctx,
			`INSERT INTO reminder_tokens (profile_id, token_hash, expires_at, created_at)
			 VALUES (?, ?, ?, ?)`,
			profileID, hash[:], expiresAt, createdAt,
		)
		if err != nil {
			// In the extremely unlikely event of a UNIQUE(token_hash) conflict,
			// retry with a new random token.
			if isUniqueConstraintError(err) {
				continue
			}
			return "", err
		}

		return token, nil
	}
}

func (s *service) ValidateAndConsume(ctx context.Context, rawToken string) (int64, error) {
	if rawToken == "" {
		return 0, ErrInvalidToken
	}

	hash := sha256.Sum256([]byte(rawToken))

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var (
		id        int64
		profileID int64
		expiresAt time.Time
		usedAt    sql.NullTime
	)

	row := tx.QueryRowContext(ctx,
		`SELECT id, profile_id, expires_at, used_at
		 FROM reminder_tokens
		 WHERE token_hash = ?`, hash[:],
	)
	if err := row.Scan(&id, &profileID, &expiresAt, &usedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidToken
		}
		return 0, err
	}

	now := time.Now().UTC()
	if !expiresAt.After(now) || usedAt.Valid {
		return 0, ErrInvalidToken
	}

	res, err := tx.ExecContext(ctx,
		`UPDATE reminder_tokens
		 SET used_at = ?
		 WHERE id = ? AND used_at IS NULL`, now, id,
	)
	if err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if affected == 0 {
		// Another concurrent request consumed it first.
		return 0, ErrInvalidToken
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return profileID, nil
}

// isUniqueConstraintError attempts to detect SQLite UNIQUE constraint errors
// without depending on driver-specific types in this core package.
func isUniqueConstraintError(err error) bool {
	// For github.com/mattn/go-sqlite3 this typically contains "UNIQUE constraint failed".
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}

package tokens

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE reminder_tokens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_id INTEGER NOT NULL,
		token_hash BLOB NOT NULL,
		expires_at DATETIME NOT NULL,
		used_at DATETIME,
		created_at DATETIME NOT NULL,
		UNIQUE(token_hash)
	)`)
	if err != nil {
		t.Fatalf("create table: %v", err)
	}
	return db
}

func TestTokenGenerateValidateConsume(t *testing.T) {
	db := newTestDB(t)
	t.Cleanup(func() { db.Close() })

	svc := NewService(db, Config{TTL: time.Hour})
	ctx := context.Background()

	token, err := svc.Generate(ctx, 1)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	if token == "" {
		t.Fatalf("expected non-empty token")
	}

	profileID, err := svc.ValidateAndConsume(ctx, token)
	if err != nil {
		t.Fatalf("ValidateAndConsume error: %v", err)
	}
	if profileID != 1 {
		t.Fatalf("expected profileID 1, got %d", profileID)
	}

	// Second use should fail.
	if _, err := svc.ValidateAndConsume(ctx, token); err == nil {
		t.Fatalf("expected error on second validation, got nil")
	}
}

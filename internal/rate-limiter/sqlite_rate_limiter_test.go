package ratelimiter

import (
	"database/sql"
	"testing"

	"golang.org/x/time/rate"
	_ "modernc.org/sqlite"
)

// compile-time check that SQLiteRateLimiter satisfies the interface.
var _ RateLimiter = (*SQLiteRateLimiter)(nil)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	// Shared-cache in-memory DB so the single connection pool sees one DB.
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func TestSQLiteRateLimiterAllowsBurstThenLimits(t *testing.T) {
	db := newTestDB(t)

	// rps=5, burst=3 -> first 3 requests allowed, 4th denied.
	rl, err := NewSQLiteRateLimiter(db, 5, 3)
	if err != nil {
		t.Fatalf("new limiter: %v", err)
	}

	for i := range 3 {
		if err := rl.Allow("alice"); err != nil {
			t.Fatalf("request %d: expected allowed, got %v", i, err)
		}
	}

	if err := rl.Allow("alice"); err != ErrRateLimited {
		t.Fatalf("4th request: expected ErrRateLimited, got %v", err)
	}
}

func TestSQLiteRateLimiterIsolatesClients(t *testing.T) {
	db := newTestDB(t)
	rl, err := NewSQLiteRateLimiter(db, 1, 1)
	if err != nil {
		t.Fatalf("new limiter: %v", err)
	}

	// Each client has its own bucket.
	if err := rl.Allow("alice"); err != nil {
		t.Fatalf("alice first request should be allowed, got %v", err)
	}
	if err := rl.Allow("bob"); err != nil {
		t.Fatalf("bob first request should be allowed, got %v", err)
	}
	// alice is now exhausted.
	if err := rl.Allow("alice"); err != ErrRateLimited {
		t.Fatalf("alice second request: expected ErrRateLimited, got %v", err)
	}
}

func TestSQLiteRateLimiterLimit(t *testing.T) {
	db := newTestDB(t)
	rl, err := NewSQLiteRateLimiter(db, 7, 10)
	if err != nil {
		t.Fatalf("new limiter: %v", err)
	}
	if rl.Limit() != rate.Limit(7) {
		t.Fatalf("expected limit 7, got %v", rl.Limit())
	}
}

package ratelimiter

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// SQLiteRateLimiter is a per-client rate limiter that persists its token-bucket
// state in a SQLite table. It mirrors the token-bucket semantics of
// InMemRateLimiter: each client's bucket refills at `rps` tokens/second up to a
// maximum of `burst` tokens, and each allowed request consumes one token.
//
// State is persisted so limits survive process restarts. Bucket updates run
// inside a single transaction and are additionally serialized with a mutex,
// since SQLite permits only one writer at a time.
type SQLiteRateLimiter struct {
	db    *sql.DB
	mu    sync.Mutex
	rps   rate.Limit // requests per second (refill rate)
	burst int        // max burst size (bucket capacity)
}

const createBucketTable = `
CREATE TABLE IF NOT EXISTS rate_buckets (
    client_id TEXT PRIMARY KEY,
    tokens    REAL NOT NULL,
    ts        REAL NOT NULL
)`

// NewSQLiteRateLimiter creates a SQLite-backed rate limiter using the given
// open *sql.DB, allowing rps requests/sec per client with the given burst size.
// It ensures the backing table exists.
func NewSQLiteRateLimiter(db *sql.DB, rps float64, burst int) (*SQLiteRateLimiter, error) {
	if _, err := db.Exec(createBucketTable); err != nil {
		return nil, fmt.Errorf("sqlite rate limiter: create table: %w", err)
	}
	return &SQLiteRateLimiter{
		db:    db,
		rps:   rate.Limit(rps),
		burst: burst,
	}, nil
}

// Allow reports whether a request from clientID is allowed right now.
// It returns nil if allowed, ErrRateLimited if the limit is exceeded, or a
// wrapped error if the database access fails. It satisfies the RateLimiter
// interface.
func (rl *SQLiteRateLimiter) Allow(clientID string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	tx, err := rl.db.Begin()
	if err != nil {
		return fmt.Errorf("sqlite rate limiter: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	now := float64(time.Now().UnixNano()) / float64(time.Second)
	burst := float64(rl.burst)

	// Load the current bucket, defaulting to a full bucket for new clients.
	var tokens, ts float64
	row := tx.QueryRow(`SELECT tokens, ts FROM rate_buckets WHERE client_id = ?`, clientID)
	switch err := row.Scan(&tokens, &ts); err {
	case nil:
		// Refill based on elapsed time, capped at burst.
		if elapsed := now - ts; elapsed > 0 {
			tokens += elapsed * float64(rl.rps)
			if tokens > burst {
				tokens = burst
			}
		}
	case sql.ErrNoRows:
		tokens = burst
	default:
		return fmt.Errorf("sqlite rate limiter: scan: %w", err)
	}

	allowed := tokens >= 1
	if allowed {
		tokens -= 1
	}

	if _, err := tx.Exec(
		`INSERT INTO rate_buckets (client_id, tokens, ts) VALUES (?, ?, ?)
         ON CONFLICT(client_id) DO UPDATE SET tokens = excluded.tokens, ts = excluded.ts`,
		clientID, tokens, now,
	); err != nil {
		return fmt.Errorf("sqlite rate limiter: upsert: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("sqlite rate limiter: commit: %w", err)
	}

	if allowed {
		return nil
	}
	return ErrRateLimited
}

// Limit returns the configured rate limit.
func (rl *SQLiteRateLimiter) Limit() rate.Limit {
	return rl.rps
}

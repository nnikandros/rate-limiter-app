package server

import (
	"net/http"
	"net/http/httptest"
	"rate-limiter-app/internal/clientdb"
	ratelimiter "rate-limiter-app/internal/rate-limiter"
	"testing"
)

// newTestServer builds a Server using the repo-backed config but without a real
// Redis connection, so only the in-memory strategy is exercised over HTTP.
func newTestServer(t *testing.T) *Server {
	t.Helper()
	repo, err := clientdb.NewMockRepository("../../mock_db.json")
	if err != nil {
		t.Fatalf("load mock repo: %v", err)
	}
	return &Server{
		cRepo:    repo,
		limiters: make(map[string]ratelimiter.RateLimiter),
	}
}

func TestLimiterForSelectsStrategy(t *testing.T) {
	s := newTestServer(t)

	// alice -> inmem
	rl, ok := s.limiterFor("alice")
	if !ok {
		t.Fatal("expected alice to be found")
	}
	if _, isInmem := rl.(*ratelimiter.InMemRateLimiter); !isInmem {
		t.Fatalf("expected alice to use InMemRateLimiter, got %T", rl)
	}

	// cached instance should be identical on second lookup
	rl2, _ := s.limiterFor("alice")
	if rl != rl2 {
		t.Fatal("expected cached limiter instance for alice")
	}

	// bob -> redis
	s.redis = nil // not used until Allow is called
	rlBob, ok := s.limiterFor("bob")
	if !ok {
		t.Fatal("expected bob to be found")
	}
	if _, isRedis := rlBob.(*ratelimiter.RedisRateLimiter); !isRedis {
		t.Fatalf("expected bob to use RedisRateLimiter, got %T", rlBob)
	}

	// unknown client
	if _, ok := s.limiterFor("nobody"); ok {
		t.Fatal("expected unknown client to be rejected")
	}
}

func TestInmemClientRateLimited(t *testing.T) {
	s := newTestServer(t)
	handler := s.RegisterRoutes()

	// alice: rps=5, burst=10 -> first 10 requests allowed, 11th limited.
	var lastTooMany bool
	for i := range 11 {
		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		req.Header.Set("Authorization", "Bearer alice")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if i < 10 && rec.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i, rec.Code)
		}
		if i == 10 {
			lastTooMany = rec.Code == http.StatusTooManyRequests
		}
	}
	if !lastTooMany {
		t.Fatal("expected 11th request to be rate limited (429)")
	}
}

func TestUnknownClientUnauthorized(t *testing.T) {
	s := newTestServer(t)
	handler := s.RegisterRoutes()

	req := httptest.NewRequest(http.MethodGet, "/foo", nil)
	req.Header.Set("Authorization", "Bearer ghost")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for unknown client, got %d", rec.Code)
	}
}

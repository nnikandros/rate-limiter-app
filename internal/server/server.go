package server

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"rate-limiter-app/internal/clientdb"
	ratelimiter "rate-limiter-app/internal/rate-limiter"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type Server struct {
	host       string
	port       string
	clientRepo clientdb.Repository

	// db is the shared SQLite handle backing the "sqlite" strategy.
	dbRateLimiter *sql.DB

	// limiters caches a per-client RateLimiter so each client keeps its own
	// token-bucket state across requests. The chosen strategy and the client's
	// rps/burst come from its ClientConfig.
	mu       sync.Mutex
	limiters map[string]ratelimiter.RateLimiter
}

func New(host, port string, clientRepo clientdb.Repository, dbRateLimiter *sql.DB) *http.Server {

	NewServer := &Server{
		clientRepo:    clientRepo,
		dbRateLimiter: dbRateLimiter,
		limiters:      make(map[string]ratelimiter.RateLimiter),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         net.JoinHostPort(host, port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

// limiterFor returns the RateLimiter for a client, constructing it on first use
// based on the client's configured strategy ("inmem" or "sqlite") and its
// rps/burst. Returns false if the client has no configuration.
func (s *Server) limiterFor(clientID string) (ratelimiter.RateLimiter, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if rl, ok := s.limiters[clientID]; ok {
		return rl, true
	}

	clientCfg, ok := s.clientRepo.GetClientConfig(clientID)
	if !ok {
		return nil, false
	}

	var rl ratelimiter.RateLimiter
	switch clientCfg.Store {
	case "sqlite":
		sqliteRl, err := ratelimiter.NewSQLiteRateLimiter(s.dbRateLimiter, clientCfg.Rps, clientCfg.Burst)
		if err != nil {
			log.Printf("sqlite limiter for %q: %v", clientID, err)
			return nil, false
		}
		rl = sqliteRl
	default: // "inmem" and anything else fall back to in-memory
		rl = ratelimiter.NewInMemRateLimiter(clientCfg.Rps, clientCfg.Burst)
	}

	s.limiters[clientID] = rl
	return rl, true
}

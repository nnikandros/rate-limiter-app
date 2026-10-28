package server

import (
	"log"
	"net"
	"net/http"
	"rate-limiter-app/internal/clientdb"
	"rate-limiter-app/internal/config"
	ratelimiter "rate-limiter-app/internal/rate-limiter"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Server struct {
	cRepo *clientdb.MockRepository

	// redis is the shared client used to back the "redis" strategy.
	redis redis.Cmdable

	// limiters caches a per-client RateLimiter so each client keeps its own
	// token-bucket state across requests. The chosen strategy and the client's
	// rps/burst come from its ClientConfig.
	mu       sync.Mutex
	limiters map[string]ratelimiter.RateLimiter
}

func New(cfg config.Config) *http.Server {
	cRepo, err := clientdb.NewMockRepository(cfg.MockDbPath)
	if err != nil {
		log.Fatalf("load client config: %v", err)
	}

	opts, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		log.Fatalf("parse redis url: %v", err)
	}

	NewServer := &Server{
		cRepo:    cRepo,
		redis:    redis.NewClient(opts),
		limiters: make(map[string]ratelimiter.RateLimiter),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

// limiterFor returns the RateLimiter for a client, constructing it on first use
// based on the client's configured strategy ("inmem" or "redis") and its
// rps/burst. Returns false if the client has no configuration.
func (s *Server) limiterFor(clientID string) (ratelimiter.RateLimiter, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if rl, ok := s.limiters[clientID]; ok {
		return rl, true
	}

	cfg, ok := s.cRepo.GetClientConfig(clientID)
	if !ok {
		return nil, false
	}

	var rl ratelimiter.RateLimiter
	switch cfg.Store {
	case "redis":
		rl = ratelimiter.NewRedisRateLimiter(s.redis, cfg.Rps, cfg.Burst)
	default: // "inmem" and anything else fall back to in-memory
		rl = ratelimiter.NewInMemRateLimiter(cfg.Rps, cfg.Burst)
	}

	s.limiters[clientID] = rl
	return rl, true
}

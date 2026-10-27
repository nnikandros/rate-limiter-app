package server

import (
	"net"
	"net/http"
	"rate-limiter-app/internal/clients"
	"rate-limiter-app/internal/config"
	ratelimiter "rate-limiter-app/internal/rate-limiter"
	"time"
)

type Server struct {
	// host string
	// port string
	// db        database.Service
	// templates *template.Template
	// log       *slog.Logger
	cRepo       clients.Repository
	rateLimiter ratelimiter.RateLimiter
}

func New(cfg config.Config) *http.Server {
	NewServer := &Server{
		// db:        database.NewService(),
		// templates: t,
		// log:       getLogger(),
		cRepo: clients.New(cfg),
		// Allow 5 requests/sec per client, with a burst of 10.
		rateLimiter: ratelimiter.NewInMemRateLimiter(5, 10),
	}

	// Declare Server config
	server := &http.Server{
		// Addr:         fmt.Sprintf("%s:%d", host, port),
		Addr:         net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

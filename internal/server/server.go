package server

import (
	"net"
	"net/http"
	clientrepository "rate-limiter-app/internal/client-repository"
	"rate-limiter-app/internal/config"
	"time"
)

type Server struct {
	// host string
	// port string
	// db        database.Service
	// templates *template.Template
	// log       *slog.Logger
	clientRepo clientrepository.ClientRepository
}

func New(cfg config.Config) *http.Server {
	NewServer := &Server{
		// db:        database.NewService(),
		// templates: t,
		// log:       getLogger(),
		clientRepo: clientrepository.New(cfg),
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

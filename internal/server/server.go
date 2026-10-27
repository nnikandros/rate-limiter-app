package server

import (
	"net"
	"net/http"
	"rate-limiter-app/internal/config"
	"time"
)

type Server struct {
	// host string
	// port string
	// db        database.Service
	// templates *template.Template
	// log       *slog.Logger
}

func New(cfg config.Config) *http.Server {
	// port, _ := strconv.Atoi(os.Getenv("PORT"))
	// port := os.Getenv("PORT")
	// host := os.Getenv("HOST")

	NewServer := &Server{
		// db:        database.NewService(),
		// templates: t,
		// log:       getLogger(),
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

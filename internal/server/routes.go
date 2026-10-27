package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// r.With(s.CasMiddleware).Get("/", s.HelloWorldHandler)

	r.With(s.AuthEMiddleware).With(s.RateLimiterMiddleware).Get("/foo", s.successHandler)
	r.With(s.AuthEMiddleware).With(s.RateLimiterMiddleware).Get("/bar", s.successHandler)

	return r
}

func (s *Server) AuthEMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := r.Header.Get("Authorization")

		if client != "client-1" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}

	})
}

func (s *Server) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func (s *Server) successHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp := map[string]bool{"success": true}
	w.Header().Set("Content-Type", "application/json")

	b, _ := json.Marshal(jsonResp)

	_, _ = w.Write(b)

}

package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type SuccessResponse struct {
	Success bool `json:"success"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

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

	r.With(s.AuthEMiddleware).With(s.RateLimiterMiddleware).Get("/foo", s.successHandler)
	r.With(s.AuthEMiddleware).With(s.RateLimiterMiddleware).Get("/bar", s.successHandler)

	return r
}

func (s *Server) AuthEMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authZHeader := r.Header.Get("Authorization")
		client := strings.TrimPrefix(authZHeader, "Bearer ")

		// A client is authenticated if it exists in the client config.
		if _, ok := s.clientRepo.GetClientConfig(client); !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authZHeader := r.Header.Get("Authorization")
		client := strings.TrimPrefix(authZHeader, "Bearer ")

		// Select the per-client strategy (inmem or sqlite) from config.
		limiter, ok := s.limiterFor(client)
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if err := limiter.Allow(client); err != nil {
			w.WriteHeader(http.StatusTooManyRequests)
			e := ErrorResponse{Error: "rate limit exceeded"}
			b, _ := json.Marshal(e)
			_, _ = w.Write(b)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) successHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jsonResp := SuccessResponse{Success: true}
	b, _ := json.Marshal(jsonResp)

	_, _ = w.Write(b)

}

package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type SuccessResponse struct {
	Success bool
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

		if client != "client-1" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (s *Server) successHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jsonResp := map[string]bool{"success": true}
	b, _ := json.Marshal(jsonResp)

	_, _ = w.Write(b)

	// resp := SuccessResponse{Success: true}
	// err := json.NewEncoder(w).Encode(resp)

	// if err != nil {
	// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// }

}

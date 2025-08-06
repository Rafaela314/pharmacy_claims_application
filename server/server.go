package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pharmacy_claims_application/db"
	"github.com/pharmacy_claims_application/util"
)

type Server struct {
	store  db.Store
	router *http.ServeMux
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store:  store,
		router: http.NewServeMux(),
	}

	server.setupRoutes()
	return server
}

func (server *Server) setupRoutes() {
	// Health check endpoint
	server.router.HandleFunc("GET /health", server.healthCheck)
	
	// API endpoints
	server.router.HandleFunc("POST /api/v1/claim", server.createClaim)
	server.router.HandleFunc("GET /api/v1/claim/{id}", server.getClaim)
	server.router.HandleFunc("POST /api/v1/reversal", server.createReversal)
}

func (server *Server) Start(config util.Config) error {
	serverWithMiddleware := server.loggingMiddleware(server.router)
	
	srv := &http.Server{
		Addr:         config.ServerAddress,
		Handler:      serverWithMiddleware,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting server on %s", config.ServerAddress)
	return srv.ListenAndServe()
}

// Middleware for logging requests
func (server *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create a custom response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(wrapped, r)
		
		log.Printf(
			"%s %s %d %v",
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			time.Since(start),
		)
	})
}

// Custom response writer to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
} 
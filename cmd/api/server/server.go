package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	port   string
	server *http.Server
	router *chi.Mux
}

// Returns a new Server application
func New(port string) *Server {
	router := chi.NewRouter()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	return &Server{
		router: router,
		server: server,
	}
}

// Start the server listening in the configured port
func (server *Server) Listen() error {
	if err := server.server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// Gracefully shutdown the server, closing all open and idle connections before
// the shutdown is complete
func (server *Server) Shutdown(ctx context.Context) error {
	return server.server.Shutdown(ctx)
}

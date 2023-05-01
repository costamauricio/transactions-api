package server

import (
	"context"
	"net/http"

	"github.com/costamauricio/transactions-api/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

type Server struct {
	server *http.Server
	Router *chi.Mux
}

// Returns a new Server application
func New(port string) *Server {
	router := chi.NewRouter()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	return &Server{
		Router: router,
		server: server,
	}
}

// Attach the handlers and middlewares to the server mux
// Expects an ApplicationHandler with all handlers dependecies fulfilled
func (server *Server) AttachHandlers(container *handlers.ApplicationHandlers) {
	// Changes the middleware DefaultLogger to use the application custom logger
	middleware.DefaultLogger = middleware.RequestLogger(
		&middleware.DefaultLogFormatter{
			Logger: container.Logger,
		},
	)
	server.Router.Use(middleware.Logger)
	// Add cors middleware to allow requests from any origin
	server.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))
	server.Router.Use(middleware.AllowContentType("application/json"))
	server.Router.Use(render.SetContentType(render.ContentTypeJSON))

	// Attach the application routes
	server.Router.Post("/accounts", container.CreateAccountHandler)
	server.Router.Get("/accounts/{AccountId:[0-9]+}", container.GetAccountHandler)
	server.Router.Post("/transactions", container.CreateTransactionHandler)
}

// Start the server listening in the previously configured port
func (server *Server) StartListening() error {
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

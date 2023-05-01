package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/costamauricio/transactions-api/internal/data-access"
	"github.com/costamauricio/transactions-api/internal/handlers"
	"github.com/costamauricio/transactions-api/internal/server"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

// Get the environment variable value, and when it is not defined returns the default value
func getEnvWithDefault(env string, defaultValue string) string {
	value := os.Getenv(env)

	if value == "" {
		return defaultValue
	}

	return value
}

func main() {
	// Loads the .env file into the environment variables
	godotenv.Load()

	// Creates a new customized logger using logrus
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Opens the database connection
	db, err := sql.Open("sqlite3", getEnvWithDefault("DATABASE_FILE", "transactions.db"))
	if err != nil {
		logger.Errorf("Failed to connect to the database: %s", err)
		return
	}
	defer db.Close()

	// Inject the database pool to the daos and inject handlers dependencies
	container := &handlers.ApplicationHandlers{
		Logger:        logger,
		AccountDAO:    &dataAccess.Account{Pool: db},
		TransactioDAO: &dataAccess.Transaction{Pool: db},
	}

	port := getEnvWithDefault("PORT", "80")

	// Creates a new Server instance
	api := server.New(port)
	api.AttachHandlers(container)

	// Creates a channel to listen for system signals to terminate the aplication
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Infof("Starting server listening at port %s", port)
		err := api.StartListening()

		// When shutting down the server, it immediately return a http.ErrServerClosed
		// so we need to ignore this error when catching errors from the server
		if err != nil && err != http.ErrServerClosed {
			logger.Errorf("Failed to start the server: %s", err)
			// Send a SIGTERM signal to terminate the application
			terminate <- syscall.SIGTERM
		}
	}()

	// Blocks until we receive a signal to terminate the application
	<-terminate
	logger.Info("Terminating the application...")

	// Create a new context with 10 seconds expiration to use as a timeout
	// when trying to shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := api.Shutdown(ctx); err != nil {
		logger.Fatalf("Error when terminating the server: %s", err)
	}
}

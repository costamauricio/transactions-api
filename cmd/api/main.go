package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/costamauricio/transactions-api/cmd/api/server"
	"github.com/costamauricio/transactions-api/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	env := config.Load()
	api := server.New(env.Port())

	// Creates a channel to listen for system signals to terminate the aplication
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Infof("Starting server listening at port %v", env.Port())
		err := api.Listen()

		// When shutting down the server, it immediately return a http.ErrServerClosed
		// so we need to ignore this error when catching errors from the server
		if err != nil && err != http.ErrServerClosed {
			logger.Errorf("Failed to start the server: %v", err)
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
		logger.Fatalf("Error when terminating the server: %v", err)
	}
}

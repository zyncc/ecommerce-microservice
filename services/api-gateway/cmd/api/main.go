package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/zyncc/ecommerce-microservice/services/api-gateway/internal/config"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/internal/server"

	"go.uber.org/zap"
)

func main() {
	env, err := config.LoadEnv()
	if err != nil {
		panic(err)
	}
	log := config.NewLogger()
	server := server.NewServer(log, env)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done, log)

	log.Info("server running", zap.Int("port", env.Port))
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal("http server error", zap.Error(err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Info("Graceful shutdown complete.")
}

func gracefulShutdown(apiServer *http.Server, done chan bool, log *zap.Logger) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Info("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Info("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

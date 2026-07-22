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

	redis := config.ConnectRedis(env)
	server := server.NewServer(log, env, redis)

	done := make(chan bool, 1)

	go gracefulShutdown(server, done, log)

	log.Info("Server running", zap.Int("port", env.Port))
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal("http server error", zap.Error(err))
	}

	<-done
	log.Info("Graceful shutdown complete.")
}

func gracefulShutdown(apiServer *http.Server, done chan bool, log *zap.Logger) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Info("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Info("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exiting")

	done <- true
}

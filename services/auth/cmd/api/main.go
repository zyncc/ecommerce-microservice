package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/zyncc/ecommerce-microservice/services/auth/internal/config"
	"github.com/zyncc/ecommerce-microservice/services/auth/internal/server"
	"go.uber.org/zap"
)

func main() {
	env, err := config.LoadEnv()
	if err != nil {
		panic(err)
	}
	log := config.NewLogger(env.AppEnv)

	// postgres
	pool, err := config.InitDB(env.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}

	// kafka
	kafkaProducer, err := config.ConnectProducer([]string{env.KafkaBroker})
	if err != nil {
		log.Fatal("failed to connect to kafka", zap.Error(err))
	}
	defer kafkaProducer.Close()

	server := server.NewServer(log, env, pool, kafkaProducer)

	done := make(chan bool, 1)

	go gracefulShutdown(server, done, log)

	log.Info("Server running", zap.Int("port", env.Port))
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
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

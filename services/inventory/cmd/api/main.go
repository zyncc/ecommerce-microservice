package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/config"
	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/consumer"
	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/server"
	"github.com/zyncc/ecommerce-microservice/services/payment/pkg/types"
	"go.uber.org/zap"
)

func main() {
	env, err := config.LoadEnv()
	if err != nil {
		panic(err)
	}
	log := config.NewLogger(env.AppEnv)

	pool, err := config.InitDB(env.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	kafkaProducer, err := config.ConnectProducer([]string{env.KafkaBroker})
	if err != nil {
		log.Fatal("failed to connect to kafka", zap.Error(err))
	}
	defer kafkaProducer.Close()
	log.Info("Kafka Producer Running")

	apiServer := server.NewServer(log, env, pool, kafkaProducer)

	inventoryRepo := repository.NewInventoryRepository(log, pool)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	inventoryConsumer := consumer.InventoryConsumer{
		Log:           log,
		GroupID:       "inventory-service-consumer",
		InventoryRepo: inventoryRepo,
		Brokers:       []string{env.KafkaBroker},
		Topics:        []string{types.PaymentSucceededTopic},
	}

	wg.Go(func() {
		if err := inventoryConsumer.RunInventoryConsumer(ctx); !errors.Is(err, context.Canceled) {
			log.Fatal("payment consumer exited with error", zap.Error(err))
		}
	})

	wg.Add(1)
	go gracefulShutdown(ctx, apiServer, &wg, log)

	log.Info("Server running", zap.Int("port", env.Port))
	if err := apiServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("Failed to start server", zap.Error(err))
	}

	wg.Wait()
	log.Info("shutdown complete")
}

func gracefulShutdown(ctx context.Context, apiServer *http.Server, wg *sync.WaitGroup, log *zap.Logger) {
	defer wg.Done()

	<-ctx.Done()
	log.Info("shutting down gracefully, press Ctrl+C again to force")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := apiServer.Shutdown(shutdownCtx); err != nil {
		log.Error("Server forced to shutdown with error", zap.Error(err))
	}

	log.Info("Server exiting")
}

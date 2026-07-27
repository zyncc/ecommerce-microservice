package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os/signal"
	"sync"
	"syscall"

	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/client"
	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/config"
	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/consumer"
	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/controller"
	grpcserver "github.com/zyncc/ecommerce-microservice/services/inventory/internal/grpc"
	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/service"
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

	// client
	orderClient, err := client.NewOrderGRPCClient(log, env.OrderServiceURL)
	if err != nil {
		log.Fatal("failed to initialize order grpc client", zap.Error(err))
	}

	// repository
	inventoryRepo := repository.NewInventoryRepository(log, pool)

	// services
	inventoryService := service.NewInventoryService(log, inventoryRepo)

	// controllers
	inventoryController := controller.NewInventoryController(log, inventoryService)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", env.Port))
	if err != nil {
		log.Fatal("failed to start grpc server", zap.Error(err))
	}

	grpcServer := grpcserver.NewServer(log, inventoryController)

	inventoryConsumer := consumer.InventoryConsumer{
		Log:           log,
		GroupID:       "inventory-service-consumer",
		InventoryRepo: inventoryRepo,
		Brokers:       []string{env.KafkaBroker},
		Topics:        []string{types.PaymentSucceededTopic},
		OrderClient:   orderClient,
	}

	wg.Go(func() {
		if err := inventoryConsumer.RunInventoryConsumer(ctx); !errors.Is(err, context.Canceled) {
			log.Fatal("payment consumer exited with error", zap.Error(err))
		}
	})

	wg.Add(1)

	log.Info("Server running", zap.Int("port", env.Port))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}

	wg.Wait()
	log.Info("shutdown complete")
}

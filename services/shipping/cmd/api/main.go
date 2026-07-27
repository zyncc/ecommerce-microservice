package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os/signal"
	"sync"
	"syscall"

	"github.com/zyncc/ecommerce-microservice/services/shipping/internal/config"
	"github.com/zyncc/ecommerce-microservice/services/shipping/internal/consumer"
	"github.com/zyncc/ecommerce-microservice/services/shipping/internal/controller"
	grpcserver "github.com/zyncc/ecommerce-microservice/services/shipping/internal/grpc"
	"github.com/zyncc/ecommerce-microservice/services/shipping/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/shipping/internal/service"
	"github.com/zyncc/ecommerce-microservice/services/shipping/pkg/types"
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

	// repositories
	shipmentRepo := repository.NewShippingRepository(log, pool)

	// services
	shipmentService := service.NewShipmentService(log, shipmentRepo, kafkaProducer)

	// controllers
	shipmentController := controller.NewShipmentController(log, shipmentService)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	shipmentConsumer := consumer.ShipmentConsumer{
		Log:          log,
		GroupID:      "shipment-service-consumer",
		ShipmentRepo: shipmentRepo,
		Brokers:      []string{env.KafkaBroker},
		Topics:       []string{types.PaymentSucceededTopic},
	}

	wg.Go(func() {
		if err := shipmentConsumer.RunShipmentConsumer(ctx); !errors.Is(err, context.Canceled) {
			log.Fatal("order consumer exited with error", zap.Error(err))
		}
	})

	wg.Add(1)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", env.Port))
	if err != nil {
		log.Fatal("failed to start grpc server", zap.Error(err))
	}

	grpcServer := grpcserver.NewServer(log, shipmentController)

	log.Info("Server running", zap.Int("port", env.Port))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}

	wg.Wait()
	log.Info("shutdown complete")
}

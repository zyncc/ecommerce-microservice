package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"sync"
	"syscall"

	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/client"
	"github.com/zyncc/ecommerce-microservice/services/order/internal/config"
	"github.com/zyncc/ecommerce-microservice/services/order/internal/consumer"
	"github.com/zyncc/ecommerce-microservice/services/order/internal/controller"
	grpcserver "github.com/zyncc/ecommerce-microservice/services/order/internal/grpc"
	"github.com/zyncc/ecommerce-microservice/services/order/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/order/internal/service"
	"github.com/zyncc/ecommerce-microservice/services/order/pkg/types"
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
	defer pool.Close()

	// redis
	redis := config.ConnectRedis(env)
	defer redis.Close()

	// kafka
	kafkaProducer, err := config.ConnectProducer([]string{env.KafkaBroker})
	if err != nil {
		log.Fatal("failed to connect to kafka", zap.Error(err))
	}
	defer kafkaProducer.Close()
	log.Info("Kafka Producer Running")

	// clients
	grpcAuthClient, err := client.NewGRPCAuthClient(log, env.AuthServiceURL)
	if err != nil {
		log.Fatal("failed to initialize auth grcp client", zap.Error(err))
	}

	productClient, err := client.NewGRPCProductClient(log, env.ProductServiceURL)
	if err != nil {
		log.Fatal("failed to initialize product grcp client", zap.Error(err))
	}

	inventoryClient, err := client.NewGRPCInventoryClient(log, env.InventoryServiceURL)
	if err != nil {
		log.Fatal("failed to initialize inventory grcp client", zap.Error(err))
	}

	// repository
	orderRepo := repository.NewOrderRepository(log, pool)
	orderCache := repository.NewOrderCacheRepository(log, redis)

	// services
	orderService := service.NewOrderService(log, orderRepo, orderCache, grpcAuthClient, productClient, inventoryClient)

	// controllers
	orderController := controller.NewOrderController(log, orderService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", env.Port))
	if err != nil {
		log.Fatal("failed to start grpc server", zap.Error(err))
	}

	grpcServer := grpcserver.NewServer(log, orderController)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	orderConsumer := consumer.OrderConsumer{
		Log:       log,
		GroupID:   "order-service-consumer",
		OrderRepo: orderRepo,
		Brokers:   []string{env.KafkaBroker},
		Topics:    []string{types.PaymentSucceededTopic, types.ShipmentUpdatedTopic},
	}

	wg.Go(func() {
		if err := orderConsumer.RunOrderConsumer(ctx); !errors.Is(err, context.Canceled) {
			log.Fatal("order consumer exited with error", zap.Error(err))
		}
	})

	wg.Add(1)

	log.Info("Server running", zap.Int("port", env.Port))
	if err := grpcServer.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("Failed to start server", zap.Error(err))
	}

	wg.Wait()
}

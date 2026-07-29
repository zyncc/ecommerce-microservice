package main

import (
	"fmt"
	"net"

	"github.com/zyncc/ecommerce-microservice/services/product/internal/config"
	"github.com/zyncc/ecommerce-microservice/services/product/internal/controller"
	grpcserver "github.com/zyncc/ecommerce-microservice/services/product/internal/grpc"
	"github.com/zyncc/ecommerce-microservice/services/product/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/product/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
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

	// repository
	productRepo := repository.NewProductRepository(log, pool)

	// cache
	productCacheRepo := repository.NewProductCacheRepository(log, redis)

	// services
	productService := service.NewProductService(log, productRepo, productCacheRepo)

	// controllers
	productController := controller.NewProductController(log, productService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", env.Port))
	if err != nil {
		log.Fatal("failed to start grpc server", zap.Error(err))
	}

	grpcServer := grpcserver.NewServer(log, productController)

	// grpc health endpoint
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	log.Info("Server running", zap.Int("port", env.Port))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}

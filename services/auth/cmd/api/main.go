package main

import (
	"fmt"
	"net"

	"github.com/zyncc/ecommerce-microservice/services/auth/internal/config"
	"github.com/zyncc/ecommerce-microservice/services/auth/internal/controller"
	grpcserver "github.com/zyncc/ecommerce-microservice/services/auth/internal/grpc"
	"github.com/zyncc/ecommerce-microservice/services/auth/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/auth/internal/service"
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

	// kafka
	kafkaProducer, err := config.ConnectProducer([]string{env.KafkaBroker})
	if err != nil {
		log.Fatal("failed to connect to kafka", zap.Error(err))
	}
	defer kafkaProducer.Close()

	// repository
	userRepo := repository.NewUserRepository(log, pool)
	addressRepo := repository.NewAddressRepository(log, pool)

	// services
	authService := service.NewAuthService(log, userRepo, kafkaProducer, env)
	addressService := service.NewAddressService(log, addressRepo, env)

	// controllers
	authController := controller.NewAuthController(log, authService, addressService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", env.Port))
	if err != nil {
		log.Fatal("failed to start server", zap.Error(err))
	}

	grpcServer := grpcserver.NewServer(log, authController)

	log.Info("Server running", zap.Int("port", env.Port))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))

	}
}

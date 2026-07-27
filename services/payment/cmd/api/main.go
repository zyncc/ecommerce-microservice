package main

import (
	"fmt"
	"net"

	"github.com/zyncc/ecommerce-microservice/services/payment/internal/config"
	"github.com/zyncc/ecommerce-microservice/services/payment/internal/controller"
	grpcserver "github.com/zyncc/ecommerce-microservice/services/payment/internal/grpc"
	"github.com/zyncc/ecommerce-microservice/services/payment/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/payment/internal/service"
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
	log.Info("Kafka Producer Running")
	defer kafkaProducer.Close()

	// repositories
	paymentRepo := repository.NewPaymentRepository(log, pool)

	// services
	paymentService := service.NewPaymentService(log, paymentRepo, kafkaProducer)

	// controllers
	paymentController := controller.NewPaymentController(log, paymentService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", env.Port))
	if err != nil {
		log.Fatal("failed to start grpc server", zap.Error(err))
	}

	grpcServer := grpcserver.NewServer(log, paymentController)

	log.Info("Server running", zap.Int("port", env.Port))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}

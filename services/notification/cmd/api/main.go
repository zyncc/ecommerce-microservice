package main

import (
	"context"
	"errors"
	"os/signal"
	"sync"
	"syscall"

	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/client"
	"github.com/zyncc/ecommerce-microservice/services/notification/internal/config"
	"github.com/zyncc/ecommerce-microservice/services/notification/internal/consumer"
	"github.com/zyncc/ecommerce-microservice/services/notification/pkg/types"
	"go.uber.org/zap"
)

func main() {
	env, err := config.LoadEnv()
	if err != nil {
		panic(err)
	}
	log := config.NewLogger(env.AppEnv)

	kafkaProducer, err := config.ConnectProducer([]string{env.KafkaBroker})
	if err != nil {
		log.Fatal("failed to connect to kafka", zap.Error(err))
	}
	defer kafkaProducer.Close()
	log.Info("Kafka Producer Running")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	orderClient, err := client.NewOrderGRPCClient(log, env.AppEnv)
	if err != nil {
		log.Fatal("failed to initialize order grpc client", zap.Error(err))
	}

	notificationConsumer := consumer.NotificationConsumer{
		Log:         log,
		OrderClient: orderClient,
		Env:         env,
		GroupID:     "notification-service-consumer",
		Brokers:     []string{env.KafkaBroker},
		Topics:      []string{types.PaymentSucceededTopic},
	}

	wg.Go(func() {
		if err := notificationConsumer.RunNotificationConsumer(ctx); !errors.Is(err, context.Canceled) {
			log.Fatal("notification consumer exited with error", zap.Error(err))
		}
	})

	wg.Wait()
	log.Info("shutdown complete")
}

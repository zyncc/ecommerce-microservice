package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zyncc/ecommerce-microservice/services/auth/internal/config"
	"go.uber.org/zap"
)

type Server struct {
	log           *zap.Logger
	env           *config.EnvConfig
	pool          *pgxpool.Pool
	kafkaProducer sarama.SyncProducer
}

func NewServer(log *zap.Logger, env *config.EnvConfig, pool *pgxpool.Pool, kafkaProducer sarama.SyncProducer) *http.Server {
	port := env.Port
	NewServer := &Server{
		log,
		env,
		pool,
		kafkaProducer,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

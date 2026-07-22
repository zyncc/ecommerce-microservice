package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/zyncc/ecommerce-microservice/services/shipping/internal/config"
	"go.uber.org/zap"
)

type Server struct {
	log           *zap.Logger
	env           *config.EnvConfig
	pool          *pgxpool.Pool
	kafkaProducer sarama.SyncProducer
}

func NewServer(log *zap.Logger, env *config.EnvConfig, pool *pgxpool.Pool, kafkaProducer sarama.SyncProducer) *http.Server {
	NewServer := &Server{
		log,
		env,
		pool,
		kafkaProducer,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", env.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/zyncc/ecommerce-microservice/services/order/internal/config"
	"go.uber.org/zap"
)

type Server struct {
	log   *zap.Logger
	env   *config.EnvConfig
	pool  *pgxpool.Pool
	redis *redis.Client
}

func NewServer(log *zap.Logger, env *config.EnvConfig, pool *pgxpool.Pool, redis *redis.Client) *http.Server {
	NewServer := &Server{
		log,
		env,
		pool,
		redis,
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

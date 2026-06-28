package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zyncc/ecommerce-microservice/services/api-gateway/internal/config"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

type Server struct {
	port int
	log  *zap.Logger
	env  *config.EnvConfig
}

func NewServer(log *zap.Logger, env *config.EnvConfig) *http.Server {
	port := env.Port
	NewServer := &Server{
		port,
		log,
		env,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

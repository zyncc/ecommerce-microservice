package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	Port            int    `env:"PORT"`
	AppEnv          string `env:"APP_ENV"`
	DatabaseURL     string `env:"DATABASE_URL"`
	KafkaBroker     string `env:"KAFKA_BROKER"`
	OrderServiceURL string `env:"ORDER_SERVICE_URL"`
}

func LoadEnv() (*EnvConfig, error) {
	_ = godotenv.Load()

	var cfg EnvConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

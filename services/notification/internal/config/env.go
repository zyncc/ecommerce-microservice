package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	Port          int    `env:"PORT"`
	AppEnv        string `env:"APP_ENV"`
	KafkaBroker   string `env:"KAFKA_BROKER"`
	FromEmail     string `env:"FROM_EMAIL"`
	FromEmailSMTP string `env:"FROM_EMAIL_SMTP"`
	SMTPAddr      string `env:"SMTP_ADDR"`
	SMTPPort      int    `env:"SMTP_PORT"`
	SMTPPassword  string `env:"SMTP_PASSWORD"`
}

func LoadEnv() (*EnvConfig, error) {
	_ = godotenv.Load()

	var cfg EnvConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

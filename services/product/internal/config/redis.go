package config

import (
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(env *EnvConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: env.RedisURL,
	})

	return rdb
}

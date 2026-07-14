package repository

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type UserCacheRepository struct {
	logger *zap.Logger
	cache  *redis.Client
}

func NewUserCacheRepository(logger *zap.Logger, cache *redis.Client) *UserCacheRepository {
	return &UserCacheRepository{
		logger,
		cache,
	}
}

package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/zyncc/ecommerce-microservice/services/order/internal/repository/model"
	"go.uber.org/zap"
)

type OrderCacheRepository interface {
	GetOrderByID(ctx context.Context, key uuid.UUID) (model.OrderWithItems, error)
	SetOrderByID(ctx context.Context, key uuid.UUID, value model.OrderWithItems) error
}

type OrderRedisCacheRepository struct {
	log   *zap.Logger
	cache *redis.Client
}

func NewOrderCacheRepository(log *zap.Logger, cache *redis.Client) *OrderRedisCacheRepository {
	return &OrderRedisCacheRepository{
		log,
		cache,
	}
}

func (r *OrderRedisCacheRepository) GetOrderByID(ctx context.Context, key uuid.UUID) (model.OrderWithItems, error) {
	var order model.OrderWithItems

	cacheKey := fmt.Sprintf("order:%s", key.String())
	data, err := r.cache.Get(ctx, cacheKey).Bytes()
	if err != nil {
		return model.OrderWithItems{}, err
	}

	if err := json.Unmarshal(data, &order); err != nil {
		return model.OrderWithItems{}, err
	}

	return order, nil
}

func (r *OrderRedisCacheRepository) SetOrderByID(ctx context.Context, key uuid.UUID, value model.OrderWithItems) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("order:%s", key.String())

	return r.cache.Set(ctx, cacheKey, data, time.Minute*5).Err()
}

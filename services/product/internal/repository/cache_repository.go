package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/zyncc/ecommerce-microservice/services/product/internal/repository/model"
	"go.uber.org/zap"
)

type ProductCacheRepository struct {
	log   *zap.Logger
	cache *redis.Client
}

func NewProductCacheRepository(log *zap.Logger, cache *redis.Client) *ProductCacheRepository {
	return &ProductCacheRepository{
		log,
		cache,
	}
}

func (r *ProductCacheRepository) GetProductByID(ctx context.Context, id uuid.UUID) (model.Product, error) {
	key := fmt.Sprintf("product:%s", id.String())
	result, err := r.cache.Get(ctx, key).Bytes()
	if err != nil {
		return model.Product{}, err
	}

	var product model.Product

	if err := json.Unmarshal(result, &product); err != nil {
		return model.Product{}, err
	}

	return product, nil
}

func (r *ProductCacheRepository) SetProductByID(ctx context.Context, item model.Product) error {
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("product:%s", item.ID.String())
	return r.cache.Set(ctx, key, data, time.Hour).Err()
}

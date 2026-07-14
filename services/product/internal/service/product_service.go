package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/product/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/product/internal/repository/model"
	"go.uber.org/zap"
)

type ProductService struct {
	log   *zap.Logger
	repo  *repository.ProductRepository
	cache *repository.ProductCacheRepository
}

func NewProductService(log *zap.Logger, repo *repository.ProductRepository, cache *repository.ProductCacheRepository) *ProductService {
	return &ProductService{log, repo, cache}
}

func (s *ProductService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (string, error) {
	id, err := s.repo.CreateProduct(ctx, &model.CreateProductParams{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
	})
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (s *ProductService) GetAllProducts(ctx context.Context, limit, offset int) ([]*model.Product, error) {
	return s.repo.FetchAllProducts(ctx, limit, offset)
}

func (s *ProductService) GetProductByID(ctx context.Context, id uuid.UUID) (model.Product, error) {
	product, err := s.cache.GetProductByID(ctx, id)
	switch {
	case err == nil:
		return product, nil
	case errors.Is(err, redis.Nil):
		s.log.Debug("cache miss", zap.String("product_id", id.String()))
	default:
		s.log.Warn("cache error", zap.Error(err))
	}

	product, err = s.repo.GetProductByID(ctx, id)
	if err != nil {
		return model.Product{}, err
	}

	if err := s.cache.SetProductByID(ctx, product); err != nil {
		s.log.Error("failed to set cache", zap.Error(err))
	}

	return product, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteProduct(ctx, id)
}

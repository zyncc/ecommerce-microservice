package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/product/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/product/internal/repository/model"
	"go.uber.org/zap"
)

type ProductService struct {
	log  *zap.Logger
	repo *repository.ProductRepository
}

func NewProductService(log *zap.Logger, repo *repository.ProductRepository) *ProductService {
	return &ProductService{log, repo}
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
	return s.repo.GetProductByID(ctx, id)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteProduct(ctx, id)
}

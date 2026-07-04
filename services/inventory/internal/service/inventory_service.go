package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/repository/model"
	"go.uber.org/zap"
)

type InventoryService struct {
	log  *zap.Logger
	repo *repository.InventoryRepository
}

func NewProductService(log *zap.Logger, repo *repository.InventoryRepository) *InventoryService {
	return &InventoryService{log, repo}
}

func (s *InventoryService) CreateInventory(ctx context.Context, req *dto.CreateInventoryRequest) (uuid.UUID, error) {
	return s.repo.CreateInventory(ctx, &model.CreateInventoryParams{
		ProductID:  req.ProductID,
		Small:      req.Inventory.Small,
		Medium:     req.Inventory.Medium,
		Large:      req.Inventory.Large,
		ExtraLarge: req.Inventory.ExtraLarge,
	})
}

package controller

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/service"
	"go.uber.org/zap"
)

type InventoryController struct {
	log *zap.Logger
	svc *service.InventoryService
}

func NewInventoryController(log *zap.Logger, svc *service.InventoryService) *InventoryController {
	return &InventoryController{log, svc}
}

func (c *InventoryController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateInventoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := c.svc.CreateInventory(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse[uuid.UUID](w, http.StatusCreated, "created inventory", id)
}

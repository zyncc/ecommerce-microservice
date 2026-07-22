package controller

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/inventory/internal/service"
	"github.com/zyncc/ecommerce-microservice/services/inventory/pkg/types/dto"
	"go.uber.org/zap"
)

type InventoryController struct {
	log *zap.Logger
	svc *service.InventoryService
}

func NewInventoryController(log *zap.Logger, svc *service.InventoryService) *InventoryController {
	return &InventoryController{log, svc}
}

func (c *InventoryController) CreateInventory(w http.ResponseWriter, r *http.Request) {
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

	utils.SuccessResponse(w, http.StatusCreated, "created inventory", id)
}

func (c *InventoryController) GetInventoryByProductID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("productID")
	if id == "" {
		c.log.Debug("product id not provided in path value")
		utils.ErrorResponse(w, http.StatusBadRequest, "product id is required")
		return
	}

	productID, err := uuid.Parse(id)
	if err != nil {
		c.log.Debug("failed to parse UUID", zap.Error(err))
		utils.ErrorResponse(w, http.StatusBadRequest, "product id is not a valid uuid")
	}

	inventory, err := c.svc.FetchInventoryByProductID(r.Context(), productID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "failed to fetch inventory")
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "fetched inventory", inventory)
}

func (c *InventoryController) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	var req []dto.UpdateInventoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "failed to parse request body")
		return
	}
	err := c.svc.UpdateInventory(r.Context(), req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "failed to update inventory")
		return
	}

	utils.SuccessResponse[any](w, http.StatusCreated, "Updated Inventory", nil)
}

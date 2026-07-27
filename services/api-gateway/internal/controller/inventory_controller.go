package controller

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/client"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"go.uber.org/zap"
)

type InventoryController struct {
	log             *zap.Logger
	inventoryClient client.InventoryClient
}

func NewInventoryController(log *zap.Logger, inventoryClient client.InventoryClient) *InventoryController {
	return &InventoryController{
		log,
		inventoryClient,
	}
}

// FetchInventoryByProductID godoc
// @Summary Create Product
// @Description Creates a new Product
// @Tags Product
// @Accept json
// @Produce json
// @Param request body dto.CreateProductRequest true "Create Product Request"
// @Success 200 {object} uuid.UUID
// @Failure 500 {object} utils.Error
// @Router /api/v1/product [post]
func (c *InventoryController) FetchInventoryByProductID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("productID")
	if id == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "product id is required")
		return
	}

	productID, err := uuid.Parse(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "product id is not a valid uuid")
		return
	}

	inventory, err := c.inventoryClient.FetchInventoryByProductID(r.Context(), productID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "failed to fetch inventory")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Fetched Inventory", inventory)
}

package controller

import (
	"encoding/json"
	"net/http"

	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/order/internal/service"
	"go.uber.org/zap"
)

type OrderController struct {
	log *zap.Logger
	svc *service.OrderService
}

func NewOrderController(log *zap.Logger, svc *service.OrderService) *OrderController {
	return &OrderController{log, svc}
}

func (c *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.log.Debug("failed to parse request body", zap.Error(err))
		utils.ErrorResponse(w, 400, "invalid json or malformed body")
		return
	}

	id, err := c.svc.CreateOrder(r.Context(), req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "failed to create order")
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "Order Created", id)
}

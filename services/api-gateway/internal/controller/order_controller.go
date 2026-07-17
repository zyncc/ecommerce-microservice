package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/client"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/auth/pkg/types"
	"go.uber.org/zap"
)

type OrderController struct {
	log         *zap.Logger
	orderClient *client.OrderClient
}

func NewOrderController(log *zap.Logger, orderClient *client.OrderClient) *OrderController {
	return &OrderController{
		log,
		orderClient,
	}
}

// CreateOrder godoc
// @Summary Create Product
// @Description Creates a new Product
// @Tags Product
// @Accept json
// @Produce json
// @Param request body dto.CreateProductRequest true "Create Product Request"
// @Success 200 {object} uuid.UUID
// @Failure 500 {object} utils.Error
// @Router /api/v1/product [post]
func (c *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(types.SessionContextKey).(types.Session)
	if !ok {
		utils.AuthorizationErrorResponse(w)
		return
	}

	var req dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid or malformed request")
		return
	}

	if req.UserID != session.ID {
		utils.ForbiddenErrorResponse(w)
		return
	}

	orderID, err := c.orderClient.CreateOrder(r.Context(), &req)
	if err != nil {
		if httpErr, ok := errors.AsType[*utils.HTTPError](err); ok {
			c.log.Error("failed to create product", zap.Error(httpErr))
			utils.ErrorResponse(w, httpErr.Status, httpErr.Message)
			return
		}
		c.log.Error("failed to create product", zap.Error(err))
		utils.ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Order Created", orderID)
}

// FindOrderByOrderID godoc
// @Summary Create Product
// @Description Creates a new Product
// @Tags Product
// @Accept json
// @Produce json
// @Param request body dto.CreateProductRequest true "Create Product Request"
// @Success 200 {object} uuid.UUID
// @Failure 500 {object} utils.Error
// @Router /api/v1/product [post]
func (c *OrderController) FindOrderByOrderID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("orderID")
	if id == "" {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, "order id is required")
		return
	}

	orderID, err := uuid.Parse(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, "order id must be a valid uuid")
		return
	}

	order, err := c.orderClient.FindOrderByOrderID(r.Context(), orderID)
	if err != nil {
		if httpErr, ok := errors.AsType[*utils.HTTPError](err); ok {
			c.log.Error("failed to create product", zap.Error(httpErr))
			utils.ErrorResponse(w, httpErr.Status, httpErr.Message)
			return
		}
		c.log.Error("failed to create product", zap.Error(err))
		utils.ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Fetched Order", order)
}

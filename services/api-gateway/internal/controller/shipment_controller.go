package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/client"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"go.uber.org/zap"
)

type ShipmentController struct {
	log            *zap.Logger
	shipmentClient *client.ShipmentClient
}

func NewShipmentController(log *zap.Logger, shipmentClient *client.ShipmentClient) *ShipmentController {
	return &ShipmentController{
		log,
		shipmentClient,
	}
}

// ShipmentWebhook godoc
// @Summary Create Product
// @Description Creates a new Product
// @Tags Product
// @Accept json
// @Produce json
// @Param request body dto.CreateProductRequest true "Create Product Request"
// @Success 200 {object} utils.Success
// @Failure 500 {object} utils.Error
// @Router /api/v1/webhook/payment [post]
func (c *ShipmentController) ShipmentWebhook(w http.ResponseWriter, r *http.Request) {
	var req dto.ShipmentWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.log.Debug("failed to parse request body", zap.Error(err))
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid or malformed request")
		return
	}

	shipmentSignature := r.Header.Get("X-Shipment-Signature")
	if shipmentSignature == "" {
		utils.ForbiddenErrorResponse(w)
		return
	}

	err := c.shipmentClient.ShipmentWebhook(r.Context(), req, shipmentSignature)
	if err != nil {
		if httpErr, ok := errors.AsType[*utils.HTTPError](err); ok {
			c.log.Error("failed to create product", zap.Error(httpErr))
			utils.ErrorResponse(w, httpErr.Status, httpErr.Message)
			return
		}
		c.log.Error("payment webhook failed", zap.Error(err))
		utils.ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.SuccessResponse[any](w, http.StatusOK, "Webhook Successful", nil)
}

func (c *ShipmentController) GetShipmentByTrackingID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("trackingID")
	if id == "" {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, "tracking id is required")
		return
	}

	trackingID, err := uuid.Parse(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, "tracking id must be a valid uuid")
		return
	}

	shipment, err := c.shipmentClient.GetShipmentByTrackingID(r.Context(), trackingID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Fetched Shipment", shipment)
}

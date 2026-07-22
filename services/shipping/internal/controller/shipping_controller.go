package controller

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/shipping/internal/service"
	"github.com/zyncc/ecommerce-microservice/services/shipping/pkg/types/dto"
	"go.uber.org/zap"
)

type ShipmentController struct {
	log *zap.Logger
	svc *service.ShipmentService
}

func NewShipmentController(log *zap.Logger, svc *service.ShipmentService) *ShipmentController {
	return &ShipmentController{log, svc}
}

func (c *ShipmentController) ShipmentUpdateWebhook(w http.ResponseWriter, r *http.Request) {
	var req dto.ShipmentWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, "failed to parse request body")
	}

	// validate webhook signature to ensure it was sent from the shipping provider
	shippingSignature := r.Header.Get("X-Shipment-Signature")
	if shippingSignature == "" {
		utils.ForbiddenErrorResponse(w)
		return
	}

	if err := c.svc.ShipmentUpdateWebhook(r.Context(), req); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.SuccessResponse[any](w, http.StatusOK, "Webhook Processed Successfully", nil)
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

	shipment, err := c.svc.GetShipmentByTrackingID(r.Context(), trackingID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "failed to get shipment by tracking id")
		return
	}

	response := dto.ShipmentResponse{
		ID:             shipment.ID,
		OrderID:        shipment.OrderID,
		IdempotencyKey: shipment.IdempotencyKey,
		Status:         shipment.Status,
		Carrier:        shipment.Carrier,
		ShippingCost:   shipment.ShippingCost,
		TrackingNumber: shipment.TrackingNumber,
		ShippedAt:      shipment.ShippedAt,
		DeliveredAt:    shipment.DeliveredAt,
		CreatedAt:      shipment.CreatedAt,
		UpdatedAt:      shipment.UpdatedAt,
	}

	utils.SuccessResponse(w, http.StatusOK, "Fetched Shipment", response)
}

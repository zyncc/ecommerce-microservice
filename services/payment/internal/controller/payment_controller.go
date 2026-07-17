package controller

import (
	"encoding/json"
	"net/http"

	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/payment/internal/service"
	"github.com/zyncc/ecommerce-microservice/services/payment/pkg/types/dto"
	"go.uber.org/zap"
)

type PaymentController struct {
	log *zap.Logger
	svc *service.PaymentService
}

func NewPaymentController(log *zap.Logger, svc *service.PaymentService) *PaymentController {
	return &PaymentController{log, svc}
}

func (c *PaymentController) PaymentWebhook(w http.ResponseWriter, r *http.Request) {
	var req dto.PaymentWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.log.Debug("invalid request body", zap.Error(err))
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, "invalid request body")
		return
	}

	// validate payment webhook signature to ensure it was sent from the payment provider
	razorpaySignature := r.Header.Get("X-Razorpay-Signature")
	if razorpaySignature == "" {
		utils.ForbiddenErrorResponse(w)
		return
	}

	if err := c.svc.ProcessPaymentWebhook(r.Context(), req); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	utils.SuccessResponse[any](w, http.StatusOK, "Successfully Processed Webhook", nil)
}

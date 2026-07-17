package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/client"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"go.uber.org/zap"
)

type PaymentController struct {
	log           *zap.Logger
	paymentClient *client.PaymentClient
}

func NewPaymentController(log *zap.Logger, paymentClient *client.PaymentClient) *PaymentController {
	return &PaymentController{
		log,
		paymentClient,
	}
}

// PaymentWebhook godoc
// @Summary Create Product
// @Description Creates a new Product
// @Tags Product
// @Accept json
// @Produce json
// @Param request body dto.CreateProductRequest true "Create Product Request"
// @Success 200 {object} utils.Success
// @Failure 500 {object} utils.Error
// @Router /api/v1/webhook/payment [post]
func (c *PaymentController) PaymentWebhook(w http.ResponseWriter, r *http.Request) {
	var req dto.PaymentWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.log.Debug("failed to parse request body", zap.Error(err))
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid or malformed request")
		return
	}

	razorpaySignature := r.Header.Get("X-Razorpay-Signature")
	if razorpaySignature == "" {
		utils.ForbiddenErrorResponse(w)
		return
	}

	err := c.paymentClient.PaymentWebhook(r.Context(), req, razorpaySignature)
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

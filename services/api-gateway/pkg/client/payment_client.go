package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"go.uber.org/zap"
)

type PaymentClient struct {
	log           *zap.Logger
	paymentSvcURL string
	httpClient    *http.Client
}

func NewPaymentClient(log *zap.Logger, paymentSvcURL string, httpClient *http.Client) *PaymentClient {
	return &PaymentClient{
		log,
		paymentSvcURL,
		httpClient,
	}
}

func (c *PaymentClient) PaymentWebhook(ctx context.Context, req dto.PaymentWebhookRequest, razorpaySignature string) error {
	reqBody, err := json.Marshal(req)
	if err != nil {
		c.log.Error("failed to marshal json data", zap.Error(err))
		return utils.ErrSomethingWentWrong
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/v1/webhook/payment", c.paymentSvcURL), bytes.NewReader(reqBody))
	if err != nil {
		c.log.Error("failed to create http request", zap.Error(err))
		return utils.ErrSomethingWentWrong
	}

	request.Header.Set("X-Razorpay-Signature", razorpaySignature)
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		c.log.Error("failed to send http request", zap.Error(err))
		return utils.ErrSomethingWentWrong
	}
	defer resp.Body.Close()

	var body utils.Success[any]
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		c.log.Error("failed to decode response body", zap.Error(err))
		return utils.ErrSomethingWentWrong
	}

	if !body.Success {
		c.log.Error(
			"order service returned error",
			zap.Int("status", body.Code),
			zap.String("message", body.Message),
		)
		return &utils.HTTPError{
			Status:  resp.StatusCode,
			Message: body.Message,
		}
	}

	return nil
}

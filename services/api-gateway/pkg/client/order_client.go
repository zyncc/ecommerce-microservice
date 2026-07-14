package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"go.uber.org/zap"
)

type OrderClient struct {
	log         *zap.Logger
	orderSvcURL string
	httpClient  *http.Client
}

func NewOrderClient(log *zap.Logger, orderSvcURL string, httpClient *http.Client) *OrderClient {
	return &OrderClient{
		log,
		orderSvcURL,
		httpClient,
	}
}

func (c *OrderClient) CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (uuid.UUID, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		c.log.Error("failed to marshal json data", zap.Error(err))
		return uuid.Nil, utils.ErrSomethingWentWrong
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/v1/order", c.orderSvcURL), bytes.NewReader(reqBody))
	if err != nil {
		c.log.Error("failed to create http request", zap.Error(err))
		return uuid.Nil, utils.ErrSomethingWentWrong
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		c.log.Error("failed to send http request", zap.Error(err))
		return uuid.Nil, utils.ErrSomethingWentWrong
	}
	defer resp.Body.Close()

	var body utils.Success[uuid.UUID]
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		c.log.Error("failed to decode response body", zap.Error(err))
		return uuid.Nil, utils.ErrSomethingWentWrong
	}

	if !body.Success {
		c.log.Error(
			"order service returned error",
			zap.Int("status", body.Code),
			zap.String("message", body.Message),
		)
		return uuid.Nil, &utils.HTTPError{
			Status:  resp.StatusCode,
			Message: body.Message,
		}
	}

	return body.Data, nil
}

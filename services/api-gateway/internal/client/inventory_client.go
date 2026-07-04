package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/internal/config"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"go.uber.org/zap"
)

type InventoryClient struct {
	log        *zap.Logger
	env        *config.EnvConfig
	httpClient *http.Client
}

func NewInventoryClient(log *zap.Logger, env *config.EnvConfig, httpClient *http.Client) *InventoryClient {
	return &InventoryClient{
		log,
		env,
		httpClient,
	}
}

func (c *InventoryClient) CreateInventory(ctx context.Context, req *dto.CreateInventoryRequest) (uuid.UUID, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		c.log.Error("failed to marshal json data", zap.Error(err))
		return uuid.Nil, utils.ErrSomethingWentWrong
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/v1/inventory", c.env.InventoryServiceURL), bytes.NewReader(reqBody))
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
			"inventory service returned error",
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

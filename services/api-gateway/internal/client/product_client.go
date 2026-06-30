package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/internal/config"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/product/pkg/types/dto"
	"go.uber.org/zap"
)

type ProductClient struct {
	log        *zap.Logger
	env        *config.EnvConfig
	httpClient *http.Client
}

func NewProductClient(log *zap.Logger, env *config.EnvConfig, httpClient *http.Client) *ProductClient {
	return &ProductClient{
		log,
		env,
		httpClient,
	}
}

func (c *ProductClient) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (utils.Success[string], error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		c.log.Error("failed to marshal json data", zap.Error(err))
		return utils.Success[string]{}, utils.ErrSomethingWentWrong
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/v1/product", c.env.ProductServiceURL), bytes.NewReader(reqBody))
	if err != nil {
		c.log.Error("failed to create http request", zap.Error(err))
		return utils.Success[string]{}, utils.ErrSomethingWentWrong
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		c.log.Error("failed to send http request", zap.Error(err))
		return utils.Success[string]{}, utils.ErrSomethingWentWrong
	}
	defer resp.Body.Close()

	var body utils.Success[string]
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		c.log.Error("failed to decode response body", zap.Error(err))
		return utils.Success[string]{}, utils.ErrSomethingWentWrong
	}

	if !body.Success {
		c.log.Error(
			"product service returned error",
			zap.Int("status", body.Code),
			zap.String("message", body.Message),
		)
		return utils.Success[string]{}, &utils.HTTPError{
			Status:  resp.StatusCode,
			Message: body.Message,
		}
	}

	return body, nil
}

func (c *ProductClient) GetAllProducts(ctx context.Context, limit, offset int) (utils.Success[[]dto.Product], error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/v1/product?limit=%d&offset=%d", c.env.ProductServiceURL, limit, offset), nil)
	if err != nil {
		c.log.Error("failed to create http request", zap.Error(err))
		return utils.Success[[]dto.Product]{}, utils.ErrSomethingWentWrong
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		c.log.Error("failed to send http request", zap.Error(err))
		return utils.Success[[]dto.Product]{}, utils.ErrSomethingWentWrong
	}
	defer resp.Body.Close()

	var body utils.Success[[]dto.Product]
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		c.log.Error("failed to decode response body", zap.Error(err))
		return utils.Success[[]dto.Product]{}, utils.ErrSomethingWentWrong
	}

	if !body.Success {
		c.log.Error(
			"product service returned error",
			zap.Int("status", body.Code),
			zap.String("message", body.Message),
		)
		return utils.Success[[]dto.Product]{}, &utils.HTTPError{
			Status:  resp.StatusCode,
			Message: body.Message,
		}
	}

	return body, nil
}

func (c *ProductClient) GetProductByID(ctx context.Context, id uuid.UUID) (utils.Success[dto.Product], error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/v1/product/%s", c.env.ProductServiceURL, id.String()), nil)
	if err != nil {
		c.log.Error("failed to create http request", zap.Error(err))
		return utils.Success[dto.Product]{}, utils.ErrSomethingWentWrong
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		c.log.Error("failed to send http request", zap.Error(err))
		return utils.Success[dto.Product]{}, utils.ErrSomethingWentWrong
	}
	defer resp.Body.Close()

	var body utils.Success[dto.Product]
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		c.log.Error("failed to decode response body", zap.Error(err))
		return utils.Success[dto.Product]{}, utils.ErrSomethingWentWrong
	}

	if !body.Success {
		c.log.Error(
			"product service returned error",
			zap.Int("status", body.Code),
			zap.String("message", body.Message),
		)
		return utils.Success[dto.Product]{}, &utils.HTTPError{
			Status:  resp.StatusCode,
			Message: body.Message,
		}
	}
	return body, nil
}

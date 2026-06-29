package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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

func (c *ProductClient) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*string, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		c.log.Error("failed to marshal json data", zap.Error(err))
		return nil, utils.ErrSomethingWentWrong
	}

	request, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/v1/product", c.env.ProductServiceURL), bytes.NewReader(reqBody))
	if err != nil {
		c.log.Error("failed to create http request", zap.Error(err))
		return nil, utils.ErrSomethingWentWrong
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		c.log.Error("failed to send http request", zap.Error(err))
		return nil, utils.ErrSomethingWentWrong
	}
	defer response.Body.Close()

	var respBody utils.Success[string]
	err = json.NewDecoder(response.Body).Decode(&respBody)
	if err != nil {
		c.log.Error("failed to decode response body", zap.Error(err))
		return nil, utils.ErrSomethingWentWrong
	}

	if !respBody.Success {
		c.log.Error("product service service returned error", zap.String("message", respBody.Message))
		return nil, errors.New(respBody.Message)
	}

	return respBody.Data, nil
}

func (c *ProductClient) GetAllProducts(ctx context.Context) ([]dto.Product, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/v1/product", c.env.ProductServiceURL), nil)
	if err != nil {
		c.log.Error("failed to create http request", zap.Error(err))
		return nil, utils.ErrSomethingWentWrong
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		c.log.Error("failed to send http request", zap.Error(err))
		return nil, utils.ErrSomethingWentWrong
	}
	defer response.Body.Close()

	var respBody utils.Success[[]dto.Product]
	err = json.NewDecoder(response.Body).Decode(&respBody)
	if err != nil {
		c.log.Error("failed to decode response body", zap.Error(err))
		return nil, utils.ErrSomethingWentWrong
	}

	if !respBody.Success {
		c.log.Error("product service service returned error", zap.String("message", respBody.Message))
		return nil, errors.New(respBody.Message)
	}

	return *respBody.Data, nil
}

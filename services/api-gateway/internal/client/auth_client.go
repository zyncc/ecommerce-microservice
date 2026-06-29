package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/zyncc/ecommerce-microservice/services/api-gateway/internal/config"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"

	"github.com/zyncc/ecommerce-microservice/services/auth/pkg/types"
	"go.uber.org/zap"
)

type AuthClient struct {
	log        *zap.Logger
	env        *config.EnvConfig
	httpClient *http.Client
}

func NewAuthClient(log *zap.Logger, env *config.EnvConfig, httpClient *http.Client) *AuthClient {
	return &AuthClient{
		log,
		env,
		httpClient,
	}
}

func (c *AuthClient) SignUp(ctx context.Context, req *dto.SignUpRequest) (*utils.Success[string], error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		c.log.Error("failed to marshal json data", zap.Error(err))
		return nil, utils.ErrSomethingWentWrong
	}

	request, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/v1/signup", c.env.AuthServiceURL), bytes.NewReader(reqBody))
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

	respBody := utils.Success[string]{}
	err = json.NewDecoder(response.Body).Decode(&respBody)
	if err != nil {
		c.log.Error("failed to decode response body", zap.Error(err))
		return nil, utils.ErrSomethingWentWrong
	}

	if !respBody.Success {
		c.log.Error("auth service returned error", zap.String("message", respBody.Message))
		return nil, errors.New(respBody.Message)
	}

	return &respBody, nil
}

func (c *AuthClient) SignIn(ctx context.Context, req *dto.SignInRequest) (*utils.Success[dto.SignInResponse], error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		c.log.Error("failed to marshal json data", zap.Error(err))
		return nil, errors.New("failed to parse json body")
	}

	request, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/v1/signin", c.env.AuthServiceURL), bytes.NewReader(reqBody))
	if err != nil {
		c.log.Error("failed to create http request", zap.Error(err))
		return nil, errors.New("failed to send http request")
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		c.log.Error("failed to send http request", zap.Error(err))
		return nil, errors.New("auth service error: failed to send request")
	}
	defer response.Body.Close()

	respBody := utils.Success[dto.SignInResponse]{}
	err = json.NewDecoder(response.Body).Decode(&respBody)
	if err != nil {
		c.log.Error("failed to decode response body", zap.Error(err))
		return nil, errors.New("auth service error: failed to parse response body")
	}

	if !respBody.Success {
		c.log.Error("auth service returned error", zap.String("message", respBody.Message))
		return nil, errors.New(respBody.Message)
	}

	return &respBody, nil
}

func (c *AuthClient) GetSession(ctx context.Context, r *http.Request) (*utils.Success[types.Session], error) {
	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/v1/session", c.env.AuthServiceURL), nil)
	if err != nil {
		c.log.Error("failed to create http request", zap.Error(err))
		return nil, utils.ErrSomethingWentWrong
	}

	tokenString, err := utils.ExtractAuthHeader(r)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	response, err := c.httpClient.Do(request)
	if err != nil {
		c.log.Error("failed to send http request", zap.Error(err))
		return nil, utils.ErrSomethingWentWrong
	}
	defer response.Body.Close()

	var respBody utils.Success[types.Session]
	if err = json.NewDecoder(response.Body).Decode(&respBody); err != nil {
		c.log.Error("failed to decode response body", zap.Error(err))
		return nil, utils.ErrSomethingWentWrong
	}

	if !respBody.Success {
		c.log.Error("auth service returned error", zap.String("message", respBody.Message))
		return nil, errors.New(respBody.Message)
	}

	return &respBody, nil
}

func (c *AuthClient) RefreshToken(ctx context.Context, r *http.Request) (*utils.Success[string], error) {
	request, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/v1/refresh", c.env.AuthServiceURL), nil)
	if err != nil {
		c.log.Error("failed to create http request", zap.Error(err))
		return nil, utils.ErrSomethingWentWrong
	}

	for _, cookie := range r.Cookies() {
		request.AddCookie(cookie)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		c.log.Error("failed to send http request", zap.Error(err))
		return nil, utils.ErrSomethingWentWrong
	}
	defer response.Body.Close()

	var respBody utils.Success[string]
	if err = json.NewDecoder(response.Body).Decode(&respBody); err != nil {
		c.log.Error("failed to decode response body", zap.Error(err))
		return nil, utils.ErrSomethingWentWrong
	}

	if !respBody.Success {
		c.log.Error("auth service returned error", zap.String("message", respBody.Message))
		return nil, errors.New(respBody.Message)
	}

	return &respBody, nil
}

package client

import (
	"bytes"
	"context"
	"encoding/json"
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

func (c *AuthClient) SignUp(ctx context.Context, req *dto.SignUpRequest) (string, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		c.log.Error("failed to marshal json data", zap.Error(err))
		return "", utils.ErrSomethingWentWrong
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/v1/signup", c.env.AuthServiceURL),
		bytes.NewReader(reqBody),
	)
	if err != nil {
		c.log.Error("failed to create http request", zap.Error(err))
		return "", utils.ErrSomethingWentWrong
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		c.log.Error("failed to send http request", zap.Error(err))
		return "", utils.ErrSomethingWentWrong
	}
	defer resp.Body.Close()

	var body utils.Success[string]
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		c.log.Error("failed to decode response body", zap.Error(err))
		return "", utils.ErrSomethingWentWrong
	}

	if !body.Success {
		c.log.Error(
			"auth service returned server error",
			zap.Int("status", body.Code),
			zap.String("message", body.Message),
		)
		return "", &utils.HTTPError{
			Status:  body.Code,
			Message: body.Message,
		}
	}

	return body.Data, nil
}

func (c *AuthClient) SignIn(ctx context.Context, req *dto.SignInRequest) (dto.SignInResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		c.log.Error("failed to marshal json data", zap.Error(err))
		return dto.SignInResponse{}, utils.ErrSomethingWentWrong
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/v1/signin", c.env.AuthServiceURL),
		bytes.NewReader(reqBody),
	)
	if err != nil {
		c.log.Error("failed to create http request", zap.Error(err))
		return dto.SignInResponse{}, utils.ErrSomethingWentWrong
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		c.log.Error("failed to send http request", zap.Error(err))
		return dto.SignInResponse{}, utils.ErrSomethingWentWrong
	}
	defer response.Body.Close()

	var body utils.Success[dto.SignInResponse]
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		c.log.Error("failed to decode response body", zap.Error(err))
		return dto.SignInResponse{}, utils.ErrSomethingWentWrong
	}

	if !body.Success {
		c.log.Error(
			"auth service returned server error",
			zap.Int("status", body.Code),
			zap.String("message", body.Message),
		)
		return dto.SignInResponse{}, &utils.HTTPError{
			Status:  body.Code,
			Message: body.Message,
		}
	}

	return body.Data, nil
}

func (c *AuthClient) GetSession(ctx context.Context, r *http.Request) (types.Session, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/api/v1/session", c.env.AuthServiceURL),
		nil,
	)
	if err != nil {
		c.log.Error("failed to create http request", zap.Error(err))
		return types.Session{}, utils.ErrSomethingWentWrong
	}

	tokenString, err := utils.ExtractAuthHeader(r)
	if err != nil {
		return types.Session{}, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	response, err := c.httpClient.Do(request)
	if err != nil {
		c.log.Error("failed to send http request", zap.Error(err))
		return types.Session{}, utils.ErrSomethingWentWrong
	}
	defer response.Body.Close()

	var body utils.Success[types.Session]
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		c.log.Error("failed to decode response body", zap.Error(err))
		return types.Session{}, utils.ErrSomethingWentWrong
	}

	if !body.Success {
		c.log.Error(
			"auth service returned server error",
			zap.Int("status", body.Code),
			zap.String("message", body.Message),
		)
		return types.Session{}, &utils.HTTPError{
			Status:  body.Code,
			Message: body.Message,
		}
	}

	return body.Data, nil
}

func (c *AuthClient) RefreshToken(ctx context.Context, r *http.Request) (string, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/v1/refresh", c.env.AuthServiceURL),
		nil,
	)
	if err != nil {
		c.log.Error("failed to create http request", zap.Error(err))
		return "", utils.ErrSomethingWentWrong
	}

	for _, cookie := range r.Cookies() {
		request.AddCookie(cookie)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		c.log.Error("failed to send http request", zap.Error(err))
		return "", utils.ErrSomethingWentWrong
	}
	defer response.Body.Close()

	var body utils.Success[string]
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		c.log.Error("failed to decode response body", zap.Error(err))
		return "", utils.ErrSomethingWentWrong
	}

	if !body.Success {
		c.log.Error(
			"auth service returned server error",
			zap.Int("status", body.Code),
			zap.String("message", body.Message),
		)
		return "", &utils.HTTPError{
			Status:  body.Code,
			Message: body.Message,
		}
	}

	return body.Data, nil
}

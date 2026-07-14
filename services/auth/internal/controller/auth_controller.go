package controller

import (
	"encoding/json"
	"net/http"

	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/auth/internal/service"
	authUtils "github.com/zyncc/ecommerce-microservice/services/auth/pkg/utils"
	"go.uber.org/zap"
)

type AuthController struct {
	logger *zap.Logger
	svc    *service.AuthService
}

func NewAuthController(logger *zap.Logger, svc *service.AuthService) *AuthController {
	return &AuthController{
		logger,
		svc,
	}
}

func (c *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	req := dto.SignUpRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid or malformed json")
		return
	}

	id, err := c.svc.SignUp(r.Context(), req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(w, http.StatusCreated, "User signed up", &id)
}

func (c *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
	req := dto.SignInRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid or malformed json")
		return
	}

	accessToken, refreshToken, err := c.svc.SignIn(r.Context(), req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "Signed in", &dto.SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (c *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "refresh token not found")
		return
	}

	refreshToken := refreshTokenCookie.Value
	accessToken, err := c.svc.RefreshToken(r.Context(), refreshToken)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Token refreshed", &accessToken)
}

func (c *AuthController) GetSession(w http.ResponseWriter, r *http.Request) {
	session, err := authUtils.GetSession(r)
	if err != nil {
		utils.AuthorizationErrorResponse(w)
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Session found", &session)
}

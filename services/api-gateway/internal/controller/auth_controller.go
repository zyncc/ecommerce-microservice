package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zyncc/ecommerce-microservice/services/api-gateway/internal/client"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"

	"go.uber.org/zap"
)

type AuthController struct {
	log        *zap.Logger
	authClient *client.AuthClient
}

func NewAuthController(log *zap.Logger, authClient *client.AuthClient) *AuthController {
	return &AuthController{
		log,
		authClient,
	}
}

func (c *AuthController) SignUp(w http.ResponseWriter, r *http.Request) {
	signUpReq := dto.SignUpRequest{}
	if err := json.NewDecoder(r.Body).Decode(&signUpReq); err != nil {
		c.log.Error("failed to decode request body", zap.Error(err))
		utils.ErrorResponse(w, http.StatusBadRequest, "failed to parse request body")
		return
	}

	response, err := c.authClient.SignUp(r.Context(), &signUpReq)
	if err != nil {
		c.log.Error("failed to sign up", zap.Error(err))
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Signed Up", &response.Data)
}

func (c *AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
	signUpReq := dto.SignInRequest{}
	if err := json.NewDecoder(r.Body).Decode(&signUpReq); err != nil {
		c.log.Error("failed to decode request body", zap.Error(err))
		utils.ErrorResponse(w, http.StatusBadRequest, "failed to parse request body")
		return
	}

	response, err := c.authClient.SignIn(r.Context(), &signUpReq)
	if err != nil {
		c.log.Error("failed to sign in", zap.Error(err))
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    *response.Data,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	utils.SuccessResponse[any](w, http.StatusOK, "Signed in", nil)
}

func (c *AuthController) SignOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})

	utils.SuccessResponse[any](w, http.StatusOK, "Signed out", nil)
}

func (c *AuthController) GetSession(w http.ResponseWriter, r *http.Request) {
	sessionResonse, err := c.authClient.GetSession(r.Context(), r)
	if err != nil {
		c.log.Error("failed to get session", zap.Error(err))
		utils.ErrorResponse(w, sessionResonse.Code, sessionResonse.Message)
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Session Found", sessionResonse.Data)
}

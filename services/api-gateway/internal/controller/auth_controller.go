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

// SignUp godoc
// @Summary Sign up
// @Description Register a new user account.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.SignUpRequest true "Sign up request"
// @Success 200 {object} string
// @Router /api/v1/signup [post]
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

// SignIn godoc
// @Summary Sign in
// @Description Authenticate a user and set the authentication cookie.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.SignInRequest true "Sign in request"
// @Success 200 {object} utils.SwaggerSuccessResponse
// @Router /api/v1/signin [post]
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

// SignOut godoc
// @Summary Sign out
// @Description Sign out the current user by clearing the authentication cookie.
// @Tags Auth
// @Produce json
// @Success 200 {object} utils.SwaggerSuccessResponse
// @Router /api/v1/signout [post]
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

// GetSession godoc
// @Summary Get session
// @Description Get the currently authenticated user's session.
// @Tags Auth
// @Produce json
// @Success 200 {object} types.Session
// @Router /api/v1/session [get]
func (c *AuthController) GetSession(w http.ResponseWriter, r *http.Request) {
	sessionResonse, err := c.authClient.GetSession(r.Context(), r)
	if err != nil {
		c.log.Error("failed to get session", zap.Error(err))
		utils.ErrorResponse(w, sessionResonse.Code, sessionResonse.Message)
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Session Found", sessionResonse.Data)
}

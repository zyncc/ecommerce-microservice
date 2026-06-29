package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zyncc/ecommerce-microservice/services/api-gateway/internal/client"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/auth/pkg/types"

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
// @Failure 500 {object} utils.Error
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
// @Success 200 {object} string
// @Failure 500 {object} utils.Error
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
		Name:     "access_token",
		Value:    response.Data.AccessToken,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(30 * time.Minute),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    response.Data.RefreshToken,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * 7 * time.Hour),
	})

	utils.SuccessResponse[any](w, http.StatusOK, "Signed in", nil)
}

// SignOut godoc
// @Summary Sign out
// @Description Sign out the current user by clearing the authentication cookie.
// @Tags Auth
// @Produce json
// @Success 200 {object} string
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
// @Failure 500 {object} utils.Error
// @Router /api/v1/session [get]
func (c *AuthController) GetSession(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(types.SessionContextKey).(types.Session)
	if !ok {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Session not found")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Session Found", &session)
}

// RefreshToken godoc
// @Summary Refresh Token
// @Description Refresh the access token by by generating a new one
// @Tags Auth
// @Produce json
// @Success 200 {object} string
// @Failure 500 {object} utils.Error
// @Router /api/v1/refresh [get]
func (c *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	resp, err := c.authClient.RefreshToken(r.Context(), r)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    *resp.Data,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(30 * time.Minute),
	})

	utils.SuccessResponse[any](w, http.StatusOK, "Token refreshed", nil)
}

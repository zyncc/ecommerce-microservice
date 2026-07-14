package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/client"
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
	var signUpReq dto.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&signUpReq); err != nil {
		c.log.Error("failed to decode request body", zap.Error(err))
		utils.ErrorResponse(w, http.StatusBadRequest, "failed to parse request body")
		return
	}

	if errs := signUpReq.Validate(); errs != nil {
		c.log.Debug("request validation failed", zap.Any("errors", errs))
		utils.ValidationErrorResponse(w, errs)
		return
	}

	id, err := c.authClient.SignUp(r.Context(), &signUpReq)
	if err != nil {
		if httpErr, ok := errors.AsType[*utils.HTTPError](err); ok {
			utils.ErrorResponse(w, httpErr.Status, httpErr.Message)
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Signed Up", id)
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
	var signInReq dto.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&signInReq); err != nil {
		c.log.Error("failed to decode request body", zap.Error(err))
		utils.ErrorResponse(w, http.StatusBadRequest, "failed to parse request body")
		return
	}

	if errs := signInReq.Validate(); errs != nil {
		c.log.Debug("request validation failed", zap.Any("fields", errs))
		utils.ValidationErrorResponse(w, errs)
		return
	}

	response, err := c.authClient.SignIn(r.Context(), &signInReq)
	if err != nil {
		if httpErr, ok := errors.AsType[*utils.HTTPError](err); ok {
			utils.ErrorResponse(w, httpErr.Status, httpErr.Message)
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    response.AccessToken,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(30 * time.Minute),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    response.RefreshToken,
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

	utils.SuccessResponse(w, http.StatusOK, "Session Found", session)
}

// RefreshToken godoc
// @Summary Refresh Token
// @Description Refresh the access token by by generating a new one
// @Tags Auth
// @Produce json
// @Success 200 {object} string
// @Failure 500 {object} utils.Error
// @Router /api/v1/refresh [post]
func (c *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := c.authClient.RefreshToken(r.Context(), r)
	if err != nil {
		if httpErr, ok := errors.AsType[*utils.HTTPError](err); ok {
			utils.ErrorResponse(w, httpErr.Status, httpErr.Message)
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(30 * time.Minute),
	})

	utils.SuccessResponse[any](w, http.StatusOK, "Token refreshed", nil)
}

// CreateAddress godoc
// @Summary Refresh Token
// @Description Refresh the access token by by generating a new one
// @Tags Auth
// @Produce json
// @Success 200 {object} string
// @Failure 500 {object} utils.Error
// @Router /api/v1/refresh [post]
func (c *AuthController) CreateAddress(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(types.SessionContextKey).(types.Session)
	if !ok {
		utils.AuthorizationErrorResponse(w)
		return
	}

	var req dto.CreateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.log.Error("failed to parse request body", zap.Error(err))
		utils.ErrorResponse(w, http.StatusBadRequest, "failed to parse request body")
		return
	}

	if errs := req.Validate(); errs != nil {
		c.log.Debug("request validation failed", zap.Any("fields", errs))
		utils.ValidationErrorResponse(w, errs)
		return
	}

	if session.ID != req.UserID {
		utils.ForbiddenErrorResponse(w)
		return
	}

	id, err := c.authClient.CreateAddress(r.Context(), req)
	if err != nil {
		if httpErr, ok := errors.AsType[*utils.HTTPError](err); ok {
			utils.ErrorResponse(w, httpErr.Status, httpErr.Message)
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.SuccessResponse[any](w, http.StatusOK, "Address Created", id)
}

// GetAddressByID godoc
// @Summary Refresh Token
// @Description Refresh the access token by by generating a new one
// @Tags Auth
// @Produce json
// @Success 200 {object} string
// @Failure 500 {object} utils.Error
// @Router /api/v1/refresh [post]
func (c *AuthController) GetAddressByID(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(types.SessionContextKey).(types.Session)
	if !ok {
		utils.AuthorizationErrorResponse(w)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "address id is invalid")
		return
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "address id is not a valid uuid")
		return
	}

	address, err := c.authClient.GetAddressByID(r.Context(), parsedID)
	if err != nil {
		if httpErr, ok := errors.AsType[*utils.HTTPError](err); ok {
			utils.ErrorResponse(w, httpErr.Status, httpErr.Message)
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if address.UserID != session.ID {
		utils.ForbiddenErrorResponse(w)
		return
	}

	utils.SuccessResponse[any](w, http.StatusOK, "Address Fetched", address)
}

// FetchAllAddresses godoc
// @Summary Refresh Token
// @Description Refresh the access token by by generating a new one
// @Tags Auth
// @Produce json
// @Success 200 {object} string
// @Failure 500 {object} utils.Error
// @Router /api/v1/refresh [post]
func (c *AuthController) FetchAllAddresses(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(types.SessionContextKey).(types.Session)
	if !ok {
		utils.AuthorizationErrorResponse(w)
		return
	}

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "user id is required")
		return
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "user id is not a valid uuid")
		return
	}

	if id != session.ID {
		utils.ForbiddenErrorResponse(w)
		return
	}

	addresses, err := c.authClient.GetAllAddresses(r.Context(), id)
	if err != nil {
		if httpErr, ok := errors.AsType[*utils.HTTPError](err); ok {
			utils.ErrorResponse(w, httpErr.Status, httpErr.Message)
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.SuccessResponse[any](w, http.StatusOK, "Fetched all Addresses", addresses)
}

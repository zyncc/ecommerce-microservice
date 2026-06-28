package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/auth/internal/service"
	"github.com/zyncc/ecommerce-microservice/services/auth/pkg/types"
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

	if err := req.Validate(); err != nil {
		c.logger.Error("validation failed", zap.Error(err))
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
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

	jwtToken, err := c.svc.SignIn(r.Context(), req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse[string](w, http.StatusCreated, "Signed in", &jwtToken)
}

func (c *AuthController) GetSession(w http.ResponseWriter, r *http.Request) {
	session, err := getSession(r)
	if err != nil {
		utils.AuthorizationErrorResponse(w)
		return
	}

	utils.SuccessResponse(w, http.StatusOK, "Session found", &session)
}

func getSession(r *http.Request) (types.Session, error) {
	token, err := parseJWT(r.Header.Get("Authorization"))
	if err != nil {
		return types.Session{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return types.Session{}, errors.New("invalid claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return types.Session{}, errors.New("invalid subject")
	}

	name, ok := claims["name"].(string)
	if !ok {
		return types.Session{}, errors.New("invalid name")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return types.Session{}, errors.New("invalid email")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return types.Session{}, errors.New("invalid role")
	}

	createdAtStr, ok := claims["created_at"].(string)
	if !ok {
		return types.Session{}, errors.New("invalid created_at")
	}

	id, err := uuid.Parse(sub)
	if err != nil {
		return types.Session{}, err
	}

	createdAt, err := time.Parse(time.RFC3339Nano, createdAtStr)
	if err != nil {
		return types.Session{}, err
	}

	return types.Session{
		ID:        id,
		Name:      name,
		Email:     email,
		Role:      role,
		CreatedAt: createdAt,
	}, nil
}

func parseJWT(authHeader string) (*jwt.Token, error) {
	if authHeader == "" {
		return nil, errors.New("missing authorization header")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("invalid authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return nil, errors.New("missing token")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("jwt secret not configured")
	}

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

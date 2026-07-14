package utils

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/auth/pkg/types"
)

func GetSession(r *http.Request) (types.Session, error) {
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

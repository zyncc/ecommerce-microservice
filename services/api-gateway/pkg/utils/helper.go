package utils

import (
	"net/http"
	"strings"

	"github.com/zyncc/ecommerce-microservice/services/auth/pkg/types"
)

func GetSession(r *http.Request) *types.Session {
	session, ok := r.Context().Value(types.SessionContextKey).(types.Session)
	if !ok {
		return nil
	}

	return &session
}

func ExtractAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", ErrMissingCredentials
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", ErrMissingCredentials
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return "", ErrMissingCredentials
	}

	return tokenString, nil
}

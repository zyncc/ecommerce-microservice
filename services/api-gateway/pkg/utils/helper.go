package utils

import (
	"net/http"
	"strings"

	"github.com/zyncc/ecommerce-microservice/services/auth/pkg/types"
	"google.golang.org/grpc/codes"
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

func GRPCToHTTP(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.Aborted:
		return http.StatusConflict
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}

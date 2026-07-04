package middleware

import (
	"context"
	"net/http"

	"github.com/zyncc/ecommerce-microservice/services/api-gateway/internal/client"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
	"github.com/zyncc/ecommerce-microservice/services/auth/pkg/types"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	log        *zap.Logger
	authClient *client.AuthClient
}

func NewAuthMiddleware(log *zap.Logger, authClient *client.AuthClient) *AuthMiddleware {
	return &AuthMiddleware{
		log:        log,
		authClient: authClient,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := m.authClient.GetSession(r.Context(), r)
		if err != nil {
			utils.AuthorizationErrorResponse(w)
			return
		}

		ctx := context.WithValue(r.Context(), types.SessionContextKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := m.authClient.GetSession(r.Context(), r)
		if err != nil {
			utils.AuthorizationErrorResponse(w)
			return
		}

		if session.Role != "admin" {
			utils.ForbiddenErrorResponse(w)
			return
		}

		ctx := context.WithValue(r.Context(), types.SessionContextKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

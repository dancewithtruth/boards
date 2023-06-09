package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Wave-95/boards/backend-core/internal/endpoint"
	"github.com/Wave-95/boards/backend-core/internal/jwt"
	"github.com/Wave-95/boards/backend-core/pkg/logger"
)

type UserId int

const (
	UserIdKey UserId = 0

	ErrMsgMissingToken = "Missing or invalid bearer token"
	ErrMsgInvalidToken = "Invalid token"
)

// Auth creates a middleware function that retrieves a bearer token and validates the token.
// The middleware sets the userId in the jwt payload into the request context. If the token is
// invalid, it will write an Unauthorized response.
func Auth(jwtService jwt.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := logger.FromContext(ctx)
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
				endpoint.WriteWithError(w, http.StatusUnauthorized, ErrMsgMissingToken)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			userId, err := jwtService.VerifyToken(token)
			if err != nil {
				logger.Errorf("handler: issue verifying jwt token: %w", err)
				endpoint.WriteWithError(w, http.StatusUnauthorized, ErrMsgInvalidToken)
				return
			}
			r = r.WithContext(withUser(ctx, userId))
			next.ServeHTTP(w, r)
		})
	}
}

// UserIdFromContext returns a user ID from context
func UserIdFromContext(ctx context.Context) string {
	if userId, ok := ctx.Value(UserIdKey).(string); ok {
		return userId
	}
	return ""
}

// withUser adds the userId to a context object and returns that context
func withUser(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, UserIdKey, userId)
}

package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Wave-95/boards/backend-core/internal/endpoint"
	"github.com/Wave-95/boards/backend-core/internal/jwt"
	"github.com/Wave-95/boards/backend-core/pkg/logger"
)

type UserID int

const (
	UserIDKey UserID = 0

	ErrMsgMissingToken = "Missing or invalid bearer token"
	ErrMsgInvalidToken = "Invalid token"
)

// Auth creates a middleware function that retrieves a bearer token and validates the token.
// The middleware sets the userID in the jwt payload into the request context. If the token is
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
			userID, err := jwtService.VerifyToken(token)
			if err != nil {
				logger.Errorf("handler: issue verifying jwt token: %w", err)
				endpoint.WriteWithError(w, http.StatusUnauthorized, ErrMsgInvalidToken)
				return
			}
			r = r.WithContext(withUser(ctx, userID))
			next.ServeHTTP(w, r)
		})
	}
}

// UserIDFromContext returns a user ID from context
func UserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// withUser adds the userID to a context object and returns that context
func withUser(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

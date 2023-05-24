package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Wave-95/boards/server/internal/endpoint"
	"github.com/Wave-95/boards/server/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
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
func Auth(jwtSecret string) func(next http.Handler) http.Handler {
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
			userId, err := verifyToken(token, jwtSecret)
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

// verifyToken parses and validates a jwt token. It returns the userId on a
// valid token.
func verifyToken(tokenString string, jwtSecret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return "", fmt.Errorf("Issue parsing token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := claims["userId"]
		userIdStr, ok := userId.(string)
		if !ok {
			return "", fmt.Errorf("Issue parsing userId: %w", err)
		}
		if userIdStr == "" {
			return "", fmt.Errorf("User id not set: %w", err)
		}
		return userIdStr, nil

	} else {
		return "", errors.New("Invalid token")
	}
}

// withUser adds the userId to a context object and returns that context
func withUser(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, UserIdKey, userId)
}

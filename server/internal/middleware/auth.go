package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Wave-95/boards/server/internal/endpoint"
	"github.com/golang-jwt/jwt/v5"
)

type UserId int

const (
	UserIdKey UserId = 0

	ErrMsgMissingToken = "Missing or invalid bearer token"
	ErrMsgInvalidToken = "Invalid token"
)

func Auth(jwtSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearerString := r.Header.Get("Authorization")
			token := strings.Split(bearerString, " ")[1]

			if token == "" {
				endpoint.WriteWithError(w, http.StatusUnauthorized, ErrMsgMissingToken)
			}
			//Do auth
			userId, err := verifyToken(token, jwtSecret)
			if err != nil {
				endpoint.WriteWithError(w, http.StatusUnauthorized, ErrMsgInvalidToken)
			}
			ctx := withUser(r.Context(), userId)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// UserIdFromContext returns a logger from context. If none found, instantiate a new logger
func UserIdFromContext(ctx context.Context) string {
	if userId, ok := ctx.Value(UserIdKey).(string); ok {
		return userId
	}
	return ""
}

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

func withUser(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, UserIdKey, userId)
}

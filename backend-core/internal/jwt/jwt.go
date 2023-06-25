package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Service is an interface that represents all the capabilities for the JWT service.
type Service interface {
	GenerateToken(userID string) (string, error)
	VerifyToken(token string) (string, error)
}

type service struct {
	jwtSecret  string
	expiration int
}

// New creates a service with a provided JWT secret string and expiration (hourly) number. It implements
// the JWT Service interface.
func New(jwtSecret string, expiration int) *service {
	return &service{jwtSecret, expiration}
}

// GenerateToken takes a user ID and
func (s *service) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Duration(s.expiration) * time.Hour).Unix(),
	})
	return token.SignedString([]byte(s.jwtSecret))
}

// VerifyToken parses and validates a jwt token. It returns the userID if the token is valid.
func (s *service) VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return "", fmt.Errorf("Issue parsing token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["userID"]
		userIDStr, ok := userID.(string)
		if !ok {
			return "", fmt.Errorf("Issue parsing userID: %w", err)
		}
		if userIDStr == "" {
			return "", fmt.Errorf("User id not set: %w", err)
		}
		return userIDStr, nil
	}
	return "", errors.New("Invalid token")

}

package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	GenerateToken(userId string) (string, error)
	VerifyToken(token string) (string, error)
}

type service struct {
	jwtSecret  string
	expiration int
}

func New(jwtSecret string, expiration int) *service {
	return &service{jwtSecret, expiration}
}

func (s *service) GenerateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Duration(s.expiration) * time.Hour).Unix(),
	})
	return token.SignedString([]byte(s.jwtSecret))
}

// verifyToken parses and validates a jwt token. It returns the userId on a
// valid token.
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

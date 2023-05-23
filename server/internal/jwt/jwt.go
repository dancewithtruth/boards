package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userId string) (string, error)
}

type service struct {
	jwtSecret  string
	expiration int
}

func (s *service) GenerateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Duration(s.expiration) * time.Hour).Unix(),
	})
	return token.SignedString([]byte(s.jwtSecret))
}

func New(jwtSecret string, expiration int) JWTService {
	return &service{jwtSecret, expiration}
}

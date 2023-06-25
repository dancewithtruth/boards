package test

import (
	"net/http"

	"github.com/Wave-95/boards/backend-core/internal/jwt"
)

// NewJWTService returns a JWT service initialized with a secret and 24 hr expiration.
func NewJWTService() jwt.Service {
	jwtSecret := "fake_jwt_secret"
	jwtExp := 24 // 24 hrs
	return jwt.New(jwtSecret, jwtExp)
}

// AuthHeader builds an auth header with a bearer token.
func AuthHeader(token string) http.Header {
	header := http.Header{}
	header.Add("Authorization", "Bearer "+token)
	return header
}

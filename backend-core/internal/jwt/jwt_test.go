package jwt

import (
	"fmt"
	"testing"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJwt(t *testing.T) {
	jwtSecret := "secret"
	expiration := 1
	userID := "abc123"

	t.Run("valid token", func(t *testing.T) {
		jwtService := New(jwtSecret, expiration)
		assert.NotNil(t, jwtService)

		// test valid token
		tokenString, err := jwtService.GenerateToken(userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)
		assert.True(t, isTokenValid(tokenString, jwtSecret), "expected token to be valid but wasn't")
	})

	t.Run("expired token", func(t *testing.T) {
		jwtService := New(jwtSecret, 0) //expiratin set to 0
		tokenString, err := jwtService.GenerateToken(userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		assert.False(t, isTokenValid(tokenString, jwtSecret), "expected token to be invalid but wasn't")
	})
}

func isTokenValid(tokenString, jwtSecret string) bool {
	token, _ := jwtv5.Parse(tokenString, func(token *jwtv5.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	return token.Valid
}

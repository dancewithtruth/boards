package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	t.Run("valid jwt token sets userId to request context", func(t *testing.T) {
		jwtSecret := "secret"
		userId := "abc123"
		expiration := 1
		token, err := generateTestToken(userId, expiration, jwtSecret)
		if err != nil {
			t.Fatalf("Issue generating test token: %v", err)
		}

		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if ctxUserId, ok := r.Context().Value(UserIdKey).(string); ok {
				assert.Equal(t, userId, ctxUserId, "Expected ctx user id to match jwt user id")
			}
		})
		authMiddleware := Auth(jwtSecret)
		protectedHandler := authMiddleware(testHandler)
		protectedHandler.ServeHTTP(res, req)
	})

	t.Run("invalid jwt token returns unauthorized error", func(t *testing.T) {
		jwtSecret := "secret"
		userId := "abc123"
		expiration := 0
		token, err := generateTestToken(userId, expiration, jwtSecret)
		if err != nil {
			t.Fatalf("Issue generating test token: %v", err)
		}

		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		mux := http.NewServeMux()
		authMiddleware := Auth(jwtSecret)
		handler := authMiddleware(mux)
		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode)
		assert.Contains(t, res.Body.String(), ErrMsgInvalidToken)
	})

}

func generateTestToken(userId string, expiration int, jwtSigningKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Duration(expiration) * time.Hour).Unix(),
	})
	return token.SignedString([]byte(jwtSigningKey))
}

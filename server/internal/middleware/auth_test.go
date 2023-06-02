package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wave-95/boards/server/internal/jwt"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	jwtSecret := "secret"

	t.Run("valid jwt token sets userId into request context", func(t *testing.T) {
		userId := "abc123"
		expiration := 1
		jwtService := jwt.New(jwtSecret, expiration)

		token, err := jwtService.GenerateToken(userId)
		if err != nil {
			t.Fatalf("Issue generating test token: %v", err)
		}

		res := httptest.NewRecorder()
		req := buildAuthRequest(token)

		// testHandler is used to check if a user ID was properly set on the request context
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if ctxUserId, ok := r.Context().Value(UserIdKey).(string); ok {
				assert.Equal(t, userId, ctxUserId, "Expected ctx user id to match jwt user id")
			}
		})
		authMiddleware := Auth(jwtService)
		protectedHandler := authMiddleware(testHandler)
		protectedHandler.ServeHTTP(res, req)
	})

	t.Run("invalid jwt token returns unauthorized error", func(t *testing.T) {
		jwtService := jwt.New(jwtSecret, 0)
		token, err := jwtService.GenerateToken("abc123")
		if err != nil {
			t.Fatalf("Issue generating test token: %v", err)
		}

		res := httptest.NewRecorder()
		req := buildAuthRequest(token)
		req.Header.Set("Authorization", "Bearer "+token)

		mux := http.NewServeMux()
		authMiddleware := Auth(jwt.New(jwtSecret, 1))
		handler := authMiddleware(mux)
		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode)
		assert.Contains(t, res.Body.String(), ErrMsgInvalidToken)
	})

}

func buildAuthRequest(token string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

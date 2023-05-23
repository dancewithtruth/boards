package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wave-95/boards/server/internal/api/user"
	"github.com/Wave-95/boards/server/internal/endpoint"
	"github.com/Wave-95/boards/server/internal/jwt"
	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandleLogin(t *testing.T) {
	jwtSecret := "abc123"
	jwtExpiration := 24
	validator := validator.New()
	mockRepo := user.NewMockRepository(make(map[uuid.UUID]*user.User))
	mockRepo.CreateUser(context.Background(), newTestUser())
	jwtService := jwt.New(jwtSecret, jwtExpiration)
	service := NewService(mockRepo, jwtService)
	api := NewAPI(service, validator)

	t.Run("handler returns token on valid login", func(t *testing.T) {
		payload := strings.NewReader(`{"email":"johndoe@gmail.com", "password": "password123"}`)
		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", payload)
		api.HandleLogin(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		var response LoginResponse
		json.NewDecoder(res.Body).Decode(&response)
		assert.NotEmpty(t, response.Token)
	})

	t.Run("handler returns response on invalid login", func(t *testing.T) {
		payload := strings.NewReader(`{"email":"nonexisting@gmail.com", "password": "nonexisting"}`)
		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", payload)
		api.HandleLogin(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
		var response endpoint.ErrResponse
		json.NewDecoder(res.Body).Decode(&response)
		assert.Equal(t, ErrMsgBadLogin, response.Message)
	})
}

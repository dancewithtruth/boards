package user

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wave-95/boards/server/internal/jwt"
	"github.com/Wave-95/boards/server/internal/models"
	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreateUser(t *testing.T) {
	validator := validator.New()
	mockRepo := &mockRepository{make(map[uuid.UUID]models.User)}
	userService := NewService(mockRepo, validator)
	jwtService := jwt.New("secret", 1)
	api := NewAPI(userService, jwtService, validator)

	payload := strings.NewReader(`{"name":"john doe", "email": "john@gmail.com", "password":"password123", "is_guest":false}`)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users", payload)
	api.HandleCreateUser(res, req)
	assert.Equal(t, http.StatusCreated, res.Code)
}

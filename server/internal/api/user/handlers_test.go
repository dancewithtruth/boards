package user

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreateUser(t *testing.T) {
	validator := validator.New()
	mockRepo := &mockRepository{make(map[uuid.UUID]*User)}
	service := NewService(mockRepo, validator)
	api := NewAPI(service, validator)

	payload := strings.NewReader(`{"name":"john doe", "email": "john@gmail.com", "password":"password123", "is_guest":false}`)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users", payload)
	api.HandleCreateUser(res, req)
	assert.Equal(t, http.StatusCreated, res.Code)
}

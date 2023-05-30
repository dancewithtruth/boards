package board

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wave-95/boards/server/internal/middleware"
	"github.com/Wave-95/boards/server/internal/models"
	"github.com/Wave-95/boards/server/internal/test"
	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreateBoard(t *testing.T) {
	validator := validator.New()
	testUser := test.NewUser()
	mockUsers := make(map[uuid.UUID]models.User)
	mockUsers[testUser.Id] = testUser
	mockBoardRepo := NewMockRepository(mockUsers)
	boardService := NewService(mockBoardRepo, validator)
	boardAPI := NewAPI(boardService, validator)

	userId := uuid.New().String()
	boardName := "My first board"
	payload := strings.NewReader(`{"name":"` + boardName + `"}`)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/boards", payload)
	// add a userId to request context
	ctx := context.WithValue(req.Context(), middleware.UserIdKey, userId)
	req = req.WithContext(ctx)

	boardAPI.HandleCreateBoard(res, req)

	assert.Equal(t, http.StatusCreated, res.Result().StatusCode)
	assert.Contains(t, res.Body.String(), userId)
	assert.Contains(t, res.Body.String(), boardName)
	assert.Contains(t, res.Body.String(), defaultBoardDescription)
}

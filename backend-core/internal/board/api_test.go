package board

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wave-95/boards/backend-core/internal/middleware"
	"github.com/Wave-95/boards/backend-core/internal/test"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreateBoard(t *testing.T) {
	// Setup test
	validator := validator.New()
	testUser := test.NewUser()
	mockBoardRepo := NewMockRepository()
	mockBoardRepo.AddUser(testUser)
	boardService := NewService(mockBoardRepo, validator)
	boardAPI := NewAPI(boardService, validator)

	// Set up request
	boardName := "My first board"
	payload := strings.NewReader(`{"name":"` + boardName + `"}`)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/boards", payload)

	// add a userID to request context
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, testUser.ID.String())
	req = req.WithContext(ctx)

	boardAPI.HandleCreateBoard(res, req)

	assert.Equal(t, http.StatusCreated, res.Result().StatusCode, "expected 201 status")
	assert.Contains(t, res.Body.String(), testUser.ID.String(), "expected user ID in response to be same as input")
	assert.Contains(t, res.Body.String(), boardName, "expected board name to be same as input")
	assert.Contains(t, res.Body.String(), defaultBoardDescription, "expected board description to be default description")
}

package post

import (
	"net/http"

	"github.com/Wave-95/boards/backend-core/internal/board"
	"github.com/Wave-95/boards/backend-core/internal/endpoint"
	"github.com/Wave-95/boards/backend-core/internal/middleware"
	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/pkg/logger"
)

const (
	ErrMsgInternalServer = "Internal server error"
	ErrMsgBoardNotFound  = "Board not found"
	ErrMsgInvalidBoardId = "Invalid board ID. Please pass in a boardId query param"
)

func (api *API) HandleListPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	// get userId from context
	userId := middleware.UserIdFromContext(ctx)
	boardId := r.URL.Query().Get("boardId")
	if boardId == "" {
		endpoint.WriteWithError(w, http.StatusBadRequest, ErrMsgInvalidBoardId)
		return
	}

	boardWithMembers, err := api.boardService.GetBoardWithMembers(ctx, boardId)
	if err != nil {
		logger.Errorf("handler: failed to get board with members: %v", err)
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		return
	}

	if !board.HasBoardAccess(boardWithMembers, userId) {
		endpoint.WriteWithError(w, http.StatusNotFound, ErrMsgBoardNotFound)
		return
	}

	posts, err := api.postService.ListPosts(ctx, boardId)
	if err != nil {
		logger.Errorf("handler: failed to list posts: %v", err)
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		return
	}

	endpoint.WriteWithStatus(w, http.StatusOK, struct {
		Data []models.Post `json:"data"`
	}{Data: posts})
}

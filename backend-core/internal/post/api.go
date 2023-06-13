package post

import (
	"net/http"

	"github.com/Wave-95/boards/backend-core/internal/board"
	"github.com/Wave-95/boards/backend-core/internal/endpoint"
	"github.com/Wave-95/boards/backend-core/internal/middleware"
	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/pkg/logger"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/go-chi/chi/v5"
)

const (
	ErrMsgInternalServer = "Internal server error"
	ErrMsgBoardNotFound  = "Board not found"
	ErrMsgInvalidBoardID = "Invalid board ID. Please pass in a boardID query param"
)

type API struct {
	postService  Service
	boardService board.Service
	validator    validator.Validate
}

func NewAPI(postService Service, boardService board.Service, validator validator.Validate) API {
	return API{
		postService:  postService,
		boardService: boardService,
		validator:    validator,
	}
}

func (api *API) HandleListPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	// get userID from context
	userID := middleware.UserIDFromContext(ctx)
	boardID := r.URL.Query().Get("boardID")
	if boardID == "" {
		endpoint.WriteWithError(w, http.StatusBadRequest, ErrMsgInvalidBoardID)
		return
	}

	boardWithMembers, err := api.boardService.GetBoardWithMembers(ctx, boardID)
	if err != nil {
		logger.Errorf("handler: failed to get board with members: %v", err)
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		return
	}

	if !board.UserHasAccess(boardWithMembers, userID) {
		endpoint.WriteWithError(w, http.StatusNotFound, ErrMsgBoardNotFound)
		return
	}

	posts, err := api.postService.ListPosts(ctx, boardID)
	if err != nil {
		logger.Errorf("handler: failed to list posts: %v", err)
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		return
	}

	endpoint.WriteWithStatus(w, http.StatusOK, struct {
		Data []models.Post `json:"data"`
	}{Data: posts})
}

func (api *API) RegisterHandlers(r chi.Router, authHandler func(http.Handler) http.Handler) {
	r.Route("/posts", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(authHandler)
			r.Get("/", api.HandleListPosts)
		})
	})
}

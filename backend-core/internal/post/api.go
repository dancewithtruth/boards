package post

import (
	"net/http"

	"github.com/Wave-95/boards/backend-core/internal/board"
	"github.com/Wave-95/boards/backend-core/internal/endpoint"
	"github.com/Wave-95/boards/backend-core/internal/middleware"
	"github.com/Wave-95/boards/backend-core/pkg/logger"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/go-chi/chi/v5"
)

const (
	errMsgInternalServer = "Internal server error."
	errMsgBoardNotFound  = "Board not found."
	errMsgInvalidBoardID = "Invalid board ID. Please pass in a boardID query param."
)

// API represents the struct that encapsulates all the post API dependencies.
type API struct {
	postService  Service
	boardService board.Service
	validator    validator.Validate
}

// NewAPI creates a new API struct with the provided dependencies.
func NewAPI(postService Service, boardService board.Service, validator validator.Validate) API {
	return API{
		postService:  postService,
		boardService: boardService,
		validator:    validator,
	}
}

// HandleListPostGroups is a handler for listing post groups that belong to a board. The handler will check
// if the requesting user has access to the board.
func (api *API) HandleListPostGroups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	// get userID from context
	userID := middleware.UserIDFromContext(ctx)
	boardID := r.URL.Query().Get("boardID")
	if boardID == "" {
		endpoint.WriteWithError(w, http.StatusBadRequest, errMsgInvalidBoardID)
		return
	}

	boardWithMembers, err := api.boardService.GetBoardWithMembers(ctx, boardID)
	if err != nil {
		logger.Errorf("handler: failed to get board with members: %v", err)
		endpoint.WriteWithError(w, http.StatusInternalServerError, errMsgInternalServer)
		return
	}

	if !board.UserHasAccess(boardWithMembers, userID) {
		endpoint.WriteWithError(w, http.StatusNotFound, errMsgBoardNotFound)
		return
	}

	postGroups, err := api.postService.ListPostGroups(ctx, boardID)
	if err != nil {
		logger.Errorf("handler: failed to list post groups: %v", err)
		endpoint.WriteWithError(w, http.StatusInternalServerError, errMsgInternalServer)
		return
	}

	endpoint.WriteWithStatus(w, http.StatusOK, struct {
		Result []GroupWithPostsDTO `json:"result"`
	}{Result: postGroups})
}

// RegisterHandlers registers all the post API handlers to their respective routes.
func (api *API) RegisterHandlers(r chi.Router, authHandler func(http.Handler) http.Handler) {
	r.Route("/post-groups", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(authHandler)
			r.Get("/", api.HandleListPostGroups)
		})
	})
}

package board

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Wave-95/boards/server/internal/endpoint"
	"github.com/Wave-95/boards/server/internal/middleware"
	"github.com/Wave-95/boards/server/pkg/logger"
	"github.com/go-chi/chi/v5"
)

const (
	ErrMsgInternalServer = "Internal server error"
	ErrMsgInvalidToken   = "Invalid authentication token"
	ErrMsgBoardNotFound  = "Board not found"
	ErrMsgInvalidBoardId = "Provided invalid board ID. Please ensure board ID is in UUID format"
)

// HandleCreateBoard is the handler for creating a single board. It requires a user ID
// from the request context to assign which user owns the board
func (api *API) HandleCreateBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	// decode input
	var input CreateBoardInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		endpoint.HandleDecodeErr(w, err)
		return
	}
	defer r.Body.Close()

	// validate input
	if err := api.validator.Struct(input); err != nil {
		endpoint.WriteValidationErr(w, input, err)
		return
	}

	// get userId from context
	userId := middleware.UserIdFromContext(ctx)
	if userId == "" {
		logger.Error("handler: failed to parse user ID from request context")
		endpoint.WriteWithError(w, http.StatusUnauthorized, ErrMsgInvalidToken)
		return
	}
	input.UserId = userId

	// create board
	board, err := api.boardService.CreateBoard(ctx, input)
	if err != nil {
		logger.Errorf("handler: failed to create board: %v", err)
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		return
	}
	endpoint.WriteWithStatus(w, http.StatusCreated, board)
}

// HandleGetBoard returns a single board along with a list of associated members
func (api *API) HandleGetBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	boardId := chi.URLParam(r, "boardId")
	boardWithMembers, err := api.boardService.GetBoardWithMembers(ctx, boardId)
	if err != nil {
		if errors.Is(err, ErrInvalidBoardId) {
			endpoint.WriteWithError(w, http.StatusBadRequest, ErrMsgInvalidBoardId)
			return
		}
		logger.Errorf("handler: failed to get board by board ID: %v", err)
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		return
	}

	// Check if requesting user has permission to view board
	userId := middleware.UserIdFromContext(ctx)
	if userId != boardWithMembers.UserId.String() && !hasUser(boardWithMembers.Members, userId) {
		logger.Infof("handler: user requested for a board they do not have access to : %v", err)
		endpoint.WriteWithError(w, http.StatusNotFound, ErrMsgBoardNotFound)
		return
	}

	endpoint.WriteWithStatus(w, http.StatusOK, boardWithMembers)
}

// HandleGetBoards returns a list of owned and shared boards for a given user.
// The userId from the auth jwt will be used to query the boards. Each board will
// be hydrated with associated users and invites
func (api *API) HandleGetBoards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	// get userId from context
	userId := middleware.UserIdFromContext(ctx)

	// TODO: Make concurrent
	ownedBoards, err := api.boardService.ListOwnedBoardsWithMembers(ctx, userId)
	if err != nil {
		logger.Errorf("handler: failed to get owned boards by user ID: %v", err)
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		return
	}

	sharedBoards, err := api.boardService.ListSharedBoardsWithMembers(ctx, userId)
	if err != nil {
		logger.Errorf("handler: failed to get shared boards by user ID: %v", err)
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		return
	}
	endpoint.WriteWithStatus(w, http.StatusOK, struct {
		Owned  []BoardWithMembersDTO `json:"owned"`
		Shared []BoardWithMembersDTO `json:"shared"`
	}{Owned: ownedBoards, Shared: sharedBoards})
}

// Helper functions

func hasUser(members []BoardMemberDTO, userId string) bool {
	for _, member := range members {
		if member.Id.String() == userId {
			return true
		}
	}
	return false
}

func HasBoardAccess(board BoardWithMembersDTO, userId string) bool {
	if board.UserId.String() == userId {
		return true
	}

	for _, member := range board.Members {
		if member.Id.String() == userId {
			return true
		}
	}
	return false
}

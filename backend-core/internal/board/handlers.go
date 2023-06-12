package board

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Wave-95/boards/backend-core/internal/endpoint"
	"github.com/Wave-95/boards/backend-core/internal/middleware"
	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/pkg/logger"
	"github.com/go-chi/chi/v5"
)

const (
	// ErrMsgInternalServer is an error message for notifying an internal server error.
	ErrMsgInternalServer = "Internal server error"
	// ErrMsgInvalidToken is an error message for notifying an invalid auth token.
	ErrMsgInvalidToken = "Invalid authentication token"
	// ErrMsgBoardNotFound is an error message for notifying when a board is not found.
	ErrMsgBoardNotFound = "Board not found"
	// ErrMsgInvalidBoardID is an error message for notifying an improper board ID format.
	ErrMsgInvalidBoardID = "Provided invalid board ID. Please ensure board ID is in UUID format"
)

// HandleCreateBoard is the handler for creating a single board. It requires a user ID
// from the request context to assign which user owns the board.
func (api *API) HandleCreateBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	// Decode input
	var input CreateBoardInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		endpoint.HandleDecodeErr(w, err)
		return
	}
	defer r.Body.Close()

	// Validate input
	if err := api.validator.Struct(input); err != nil {
		endpoint.WriteValidationErr(w, input, err)
		return
	}

	// Get userID from context
	userID := middleware.UserIDFromContext(ctx)
	if userID == "" {
		logger.Error("handler: failed to parse user ID from request context")
		endpoint.WriteWithError(w, http.StatusUnauthorized, ErrMsgInvalidToken)
		return
	}
	input.UserID = userID

	// Create board
	board, err := api.boardService.CreateBoard(ctx, input)
	if err != nil {
		logger.Errorf("handler: failed to create board: %v", err)
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		return
	}
	endpoint.WriteWithStatus(w, http.StatusCreated, board)
}

// HandleGetBoard returns a single board along with a list of associated members.
func (api *API) HandleGetBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	boardID := chi.URLParam(r, "boardID")
	boardWithMembers, err := api.boardService.GetBoardWithMembers(ctx, boardID)
	if err != nil {
		if errors.Is(err, ErrInvalidID) {
			endpoint.WriteWithError(w, http.StatusBadRequest, ErrMsgInvalidBoardID)
			return
		}
		logger.Errorf("handler: failed to get board by board ID: %v", err)
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		return
	}

	// Check if requesting user has permission to view board.
	userID := middleware.UserIDFromContext(ctx)
	if userID != boardWithMembers.UserID.String() && !hasUser(boardWithMembers.Members, userID) {
		logger.Infof("handler: user requested for a board they do not have access to : %v", err)
		endpoint.WriteWithError(w, http.StatusNotFound, ErrMsgBoardNotFound)
		return
	}

	endpoint.WriteWithStatus(w, http.StatusOK, boardWithMembers)
}

// HandleGetBoards returns a list of owned and shared boards for a given user.
// The userID from the auth jwt will be used to query the boards. Each board will
// be hydrated with associated users and invites.
func (api *API) HandleGetBoards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	// get userID from context
	userID := middleware.UserIDFromContext(ctx)

	// TODO: Make concurrent
	ownedBoards, err := api.boardService.ListOwnedBoardsWithMembers(ctx, userID)
	if err != nil {
		logger.Errorf("handler: failed to get owned boards by user ID: %v", err)
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		return
	}

	sharedBoards, err := api.boardService.ListSharedBoardsWithMembers(ctx, userID)
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

// HandleCreateBoardInvites is the handler for creating board invites. It takes an array of
// receiver_id's and returns a list of created board invites.
func (api *API) HandleCreateBoardInvites(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	// Decode input
	var input CreateBoardInvitesInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		endpoint.HandleDecodeErr(w, err)
		return
	}
	defer r.Body.Close()

	// Prepare input
	userID := middleware.UserIDFromContext(ctx)
	if userID == "" {
		logger.Error("handler: failed to parse user ID from request context")
		endpoint.WriteWithError(w, http.StatusUnauthorized, ErrMsgInvalidToken)
		return
	}
	boardID := chi.URLParam(r, "boardID")
	input.SenderID = userID
	input.BoardID = boardID

	// Create board invites
	invites, err := api.boardService.CreateBoardInvites(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			endpoint.WriteWithError(w, http.StatusUnauthorized, ErrUnauthorized.Error())
		case errors.Is(err, ErrInvalidID):
			endpoint.WriteWithError(w, http.StatusUnauthorized, ErrInvalidID.Error())
		default:
			logger.Errorf("handler: failed to create board invites: %v", err)
			endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		}
		return
	}
	endpoint.WriteWithStatus(w, http.StatusCreated, invites)
}

func hasUser(members []MemberDTO, userID string) bool {
	for _, member := range members {
		if member.ID.String() == userID {
			return true
		}
	}
	return false
}

// UserHasAccess checks if a certain user ID exists in a board as the board owner or board member.
func UserHasAccess(board BoardWithMembersDTO, userID string) bool {
	for _, member := range board.Members {
		if member.ID.String() == userID {
			return true
		}
	}
	return false
}

// UserIsAdmin checks if a certain user ID has admin privileges on a board.
func UserIsAdmin(board BoardWithMembersDTO, userID string) bool {
	for _, member := range board.Members {
		if member.ID.String() == userID && member.Membership.Role == string(models.RoleAdmin) {
			return true
		}
	}
	return false
}

package board

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Wave-95/boards/server/internal/endpoint"
	"github.com/Wave-95/boards/server/internal/middleware"
	"github.com/Wave-95/boards/server/pkg/logger"
	"github.com/google/uuid"
)

const (
	ErrMsgInternalServer = "Internal server error"
)

type CreateBoardInput struct {
	Name        *string `json:"name" validate:"omitempty,required,min=3,max=20"`
	Description *string `json:"description" validate:"omitempty,required,min=3,max=100"`
	UserId      string
}

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
		endpoint.HandleValidationErr(w, err)
		return
	}

	// get userId from context
	userId := middleware.UserIdFromContext(ctx)
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

type OwnedBoardWithMembersDTO struct {
	Id          uuid.UUID        `json:"id"`
	Name        *string          `json:"name"`
	Description *string          `json:"description"`
	UserId      uuid.UUID        `json:"user_id"`
	Members     []BoardMemberDTO `json:"members"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type BoardMemberDTO struct {
	Id         uuid.UUID     `json:"id"`
	Name       string        `json:"name"`
	Email      *string       `json:"email"`
	Membership MembershipDTO `json:"membership"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type MembershipDTO struct {
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"added_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (api *API) HandleGetBoards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	// get userId from context
	userId := middleware.UserIdFromContext(ctx)

	defer r.Body.Close()

	boards, err := api.boardService.ListOwnedBoardsWithMembers(ctx, userId)
	if err != nil {
		logger.Errorf("handler: failed to get boards by user ID: %v", err)
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		return
	}
	endpoint.WriteWithStatus(w, http.StatusCreated, struct {
		Boards []OwnedBoardWithMembersDTO `json:"boards"`
	}{Boards: boards})
}

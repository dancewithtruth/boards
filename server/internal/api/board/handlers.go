package board

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Wave-95/boards/server/internal/endpoint"
	"github.com/Wave-95/boards/server/internal/middleware"
	"github.com/Wave-95/boards/server/internal/models"
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

type BoardResponse struct {
	Id          uuid.UUID           `json:"id"`
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	UserId      uuid.UUID           `json:"user_id"`
	Members     []BoardUserResponse `json:"members"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
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
	endpoint.WriteWithStatus(w, http.StatusCreated, toBoardResponse(board))
}

type BoardsResponse struct {
	Boards []BoardResponse `json:"boards"`
}

func (api *API) HandleGetBoards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	// get userId from context
	userId := middleware.UserIdFromContext(ctx)

	defer r.Body.Close()

	boards, err := api.boardService.ListBoardsByUser(ctx, userId)
	if err != nil {
		logger.Errorf("handler: failed to get boards by user ID: %v", err)
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		return
	}
	endpoint.WriteWithStatus(w, http.StatusCreated, struct {
		Boards []BoardResponse `json:"boards"`
	}{toBoardsResponse(boards)})
}

func toBoardResponse(board models.Board) BoardResponse {
	return BoardResponse{
		Id:          board.Id,
		Name:        board.Name,
		Description: board.Description,
		UserId:      board.UserId,
		CreatedAt:   board.CreatedAt,
		UpdatedAt:   board.UpdatedAt,
	}
}

type BoardUserResponse struct {
	Id        uuid.UUID   `json:"id"`
	Role      string      `json:"name"`
	User      models.User `json:"user"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

func toBoardsResponse(boards []models.Board) []BoardResponse {
	boardResponses := make([]BoardResponse, len(boards))
	fmt.Println("boards", boards)
	for i, board := range boards {
		var users []BoardUserResponse
		for _, user := range board.Users {
			user := BoardUserResponse{
				Id:        user.Id,
				Role:      string(user.Role),
				User:      user.User,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			}
			users = append(users, user)
		}
		boardResponses[i] = BoardResponse{
			Id:          board.Id,
			Name:        board.Name,
			Description: board.Description,
			UserId:      board.UserId,
			Members:     users,
			CreatedAt:   board.CreatedAt,
			UpdatedAt:   board.UpdatedAt,
		}
	}
	return boardResponses
}

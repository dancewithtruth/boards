package board

import (
	"encoding/json"
	"net/http"

	"github.com/Wave-95/boards/server/internal/endpoint"
	"github.com/Wave-95/boards/server/internal/middleware"
	"github.com/Wave-95/boards/server/pkg/logger"
)

const (
	ErrMsgInternalServer = "Internal server error"
)

func (api *API) HandleCreateBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	// get userId from context
	userId := middleware.UserIdFromContext(ctx)

	// decode request
	var req CreateBoardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		endpoint.HandleDecodeErr(w, err)
		return
	}
	defer r.Body.Close()

	// validate request
	if err := api.validator.Struct(req); err != nil {
		endpoint.HandleValidationErr(w, err)
		return
	}

	// create board
	input, err := req.ToInput(userId)
	if err != nil {
		logger.Errorf("handler: failed to convert request into input: %v", err)
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		return
	}
	board, err := api.boardService.CreateBoard(ctx, input)
	if err != nil {
		logger.Errorf("handler: failed to create board: %v", err)
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		return
	}
	endpoint.WriteWithStatus(w, http.StatusCreated, board.ToDto())
}

func (api *API) HandleGetBoards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)

	// get userId from context
	userId := middleware.UserIdFromContext(ctx)

	defer r.Body.Close()

	boards, err := api.boardService.GetBoardsByUserId(ctx, userId)
	if err != nil {
		logger.Errorf("handler: failed to get boards by user ID: %v", err)
		endpoint.WriteWithError(w, http.StatusInternalServerError, ErrMsgInternalServer)
		return
	}
	endpoint.WriteWithStatus(w, http.StatusCreated, boards.ToDto())
}

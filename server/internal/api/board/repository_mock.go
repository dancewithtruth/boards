package board

import (
	"context"

	"github.com/Wave-95/boards/server/internal/models"
	"github.com/google/uuid"
)

type mockRepository struct {
	boards map[uuid.UUID]models.Board
}

func (r *mockRepository) CreateBoard(ctx context.Context, board models.Board) error {
	r.boards[board.Id] = board
	return nil
}

func (r *mockRepository) GetBoard(ctx context.Context, boardId uuid.UUID) (models.Board, error) {
	if board, ok := r.boards[boardId]; ok {
		return board, nil
	}
	return models.Board{}, ErrBoardDoesNotExist
}

func (r *mockRepository) GetBoardsByUserId(ctx context.Context, userId uuid.UUID) ([]models.Board, error) {
	boards := []models.Board{}
	for _, board := range r.boards {
		if userId == board.UserId {
			boards = append(boards, board)
		}
	}
	return boards, nil
}

func (r *mockRepository) DeleteBoard(boardId uuid.UUID) error {
	delete(r.boards, boardId)
	return nil
}

func NewMockRepository(boards map[uuid.UUID]models.Board) Repository {
	return &mockRepository{boards}
}

package board

import (
	"context"

	"github.com/google/uuid"
)

type mockRepository struct {
	boards map[uuid.UUID]*Board
}

func (r *mockRepository) CreateBoard(ctx context.Context, board *Board) error {
	r.boards[board.Id] = board
	return nil
}

func (r *mockRepository) GetBoard(ctx context.Context, boardId uuid.UUID) (*Board, error) {
	if board, ok := r.boards[boardId]; ok {
		return board, nil
	}
	return nil, ErrBoardDoesNotExist
}

func (r *mockRepository) GetBoardsByUserId(ctx context.Context, userId uuid.UUID) ([]*Board, error) {
	boards := []*Board{}
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

func NewMockRepository(boards map[uuid.UUID]*Board) Repository {
	return &mockRepository{boards}
}

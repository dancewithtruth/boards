package board

import (
	"context"

	"github.com/Wave-95/boards/server/internal/models"
	"github.com/google/uuid"
)

type mockRepository struct {
	boards           map[uuid.UUID]models.Board
	boardMemberships map[uuid.UUID]models.BoardMembership
	users            map[uuid.UUID]models.User
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

func (r *mockRepository) GetBoardAndUsers(ctx context.Context, boardId uuid.UUID) ([]BoardAndUser, error) {
	// TODO: mock this out
	return nil, nil
}

func (r *mockRepository) ListOwnedBoards(ctx context.Context, userId uuid.UUID) ([]models.Board, error) {
	list := []models.Board{}
	for _, board := range r.boards {
		list = append(list, board)
	}
	return list, nil
}

func (r *mockRepository) ListOwnedBoardAndUsers(ctx context.Context, userId uuid.UUID) ([]BoardAndUser, error) {
	// TODO: mock this out
	return nil, nil
}

func (r *mockRepository) ListSharedBoardAndUsers(ctx context.Context, userId uuid.UUID) ([]BoardAndUser, error) {
	// TODO: mock this out
	return nil, nil
}

func (r *mockRepository) CreateBoardInvites(ctx context.Context, invites []models.BoardInvite) error {
	// TODO: mock this out
	return nil
}

func (r *mockRepository) CreateMembership(ctx context.Context, membership models.BoardMembership) error {
	r.boardMemberships[membership.Id] = membership
	return nil
}

func (r *mockRepository) DeleteBoard(ctx context.Context, boardId uuid.UUID) error {
	delete(r.boards, boardId)
	return nil
}

func NewMockRepository(users map[uuid.UUID]models.User) *mockRepository {
	boards := make(map[uuid.UUID]models.Board)
	boardMemberships := make(map[uuid.UUID]models.BoardMembership)
	return &mockRepository{boards, boardMemberships, users}
}

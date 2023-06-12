package board

import (
	"context"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/google/uuid"
)

type mockRepository struct {
	boards           map[uuid.UUID]models.Board
	boardMemberships map[uuid.UUID]models.BoardMembership
	users            map[uuid.UUID]models.User
}

func NewMockRepository(boards map[uuid.UUID]models.Board) *mockRepository {
	boardMemberships := make(map[uuid.UUID]models.BoardMembership)
	users := make(map[uuid.UUID]models.User)
	return &mockRepository{boards, boardMemberships, users}
}

// AddUser is a mock specific function to help join the relation between users and boards
func (r *mockRepository) AddUser(user models.User) {
	r.users[user.ID] = user
}

func (r *mockRepository) CreateBoard(ctx context.Context, board models.Board) error {
	r.boards[board.ID] = board
	return nil
}

func (r *mockRepository) GetBoard(ctx context.Context, boardID uuid.UUID) (models.Board, error) {
	if board, ok := r.boards[boardID]; ok {
		return board, nil
	}
	return models.Board{}, ErrBoardDoesNotExist
}

func (r *mockRepository) GetBoardAndUsers(ctx context.Context, boardID uuid.UUID) ([]BoardAndUser, error) {
	// TODO: mock this out, return board, board membership, and user
	board := r.boards[boardID]
	return []BoardAndUser{{Board: &board}}, nil
}

func (r *mockRepository) ListOwnedBoards(ctx context.Context, userID uuid.UUID) ([]models.Board, error) {
	list := []models.Board{}
	for _, board := range r.boards {
		list = append(list, board)
	}
	return list, nil
}

func (r *mockRepository) ListOwnedBoardAndUsers(ctx context.Context, userID uuid.UUID) ([]BoardAndUser, error) {
	// TODO: mock this out
	return nil, nil
}

func (r *mockRepository) ListSharedBoardAndUsers(ctx context.Context, userID uuid.UUID) ([]BoardAndUser, error) {
	// TODO: mock this out
	return nil, nil
}

func (r *mockRepository) CreateBoardInvites(ctx context.Context, invites []models.BoardInvite) error {
	// TODO: mock this out
	return nil
}

func (r *mockRepository) CreateMembership(ctx context.Context, membership models.BoardMembership) error {
	r.boardMemberships[membership.ID] = membership
	return nil
}

func (r *mockRepository) DeleteBoard(ctx context.Context, boardID uuid.UUID) error {
	delete(r.boards, boardID)
	return nil
}

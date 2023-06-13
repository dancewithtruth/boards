package board

import (
	"context"
	"time"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/google/uuid"
)

type mockRepository struct {
	boards           map[uuid.UUID]models.Board
	boardMemberships map[uuid.UUID]models.BoardMembership
	users            map[uuid.UUID]models.User
}

// NewMockRepository returns a mock board repository that implements the Repository interface
func NewMockRepository() *mockRepository {
	boards := make(map[uuid.UUID]models.Board)
	boardMemberships := make(map[uuid.UUID]models.BoardMembership)
	users := make(map[uuid.UUID]models.User)
	return &mockRepository{
		boards,
		boardMemberships,
		users,
	}
}

// AddUser is a mock specific function to help join the relation between users and boards
func (r *mockRepository) AddUser(user models.User) {
	r.users[user.ID] = user
}

func (r *mockRepository) CreateBoard(ctx context.Context, board models.Board) error {
	r.boards[board.ID] = board
	now := time.Now()
	r.boardMemberships[board.ID] = models.BoardMembership{
		ID:        uuid.New(),
		BoardID:   board.ID,
		UserID:    board.UserID,
		Role:      models.RoleAdmin,
		CreatedAt: now,
		UpdatedAt: now,
	}
	r.users[board.UserID] = models.User{
		ID: board.UserID,
	}
	return nil
}

func (r *mockRepository) GetBoard(ctx context.Context, boardID uuid.UUID) (models.Board, error) {
	if board, ok := r.boards[boardID]; ok {
		return board, nil
	}
	return models.Board{}, ErrBoardDoesNotExist
}

func (r *mockRepository) GetBoardAndUsers(ctx context.Context, boardID uuid.UUID) ([]BoardAndUser, error) {
	// TODO: mock this out, this only returns 1 row every time...
	board := r.boards[boardID]
	boardMembership := r.boardMemberships[boardID]
	user := r.users[board.UserID]
	return []BoardAndUser{{Board: &board, BoardMembership: &boardMembership, User: &user}}, nil
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

func (r *mockRepository) ListBoardInvitesFilterStatus(ctx context.Context, boardID uuid.UUID, status models.BoardInviteStatus) ([]models.BoardInvite, error) {
	// TODO: mock this out
	return nil, nil
}

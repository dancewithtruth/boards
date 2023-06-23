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
	invites          map[uuid.UUID]models.Invite
}

// NewMockRepository returns a mock board repository that implements the Repository interface.
func NewMockRepository() *mockRepository {
	boards := make(map[uuid.UUID]models.Board)
	boardMemberships := make(map[uuid.UUID]models.BoardMembership)
	users := make(map[uuid.UUID]models.User)
	invites := make(map[uuid.UUID]models.Invite)
	return &mockRepository{
		boards,
		boardMemberships,
		users,
		invites,
	}
}

// AddUser is a mock specific function to help join the relation between users and boards.
func (r *mockRepository) AddUser(user models.User) {
	r.users[user.ID] = user
}

// CreateBoard creates a mock board
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

// GetBoard returns a single mock board.
func (r *mockRepository) GetBoard(ctx context.Context, boardID uuid.UUID) (models.Board, error) {
	if board, ok := r.boards[boardID]; ok {
		return board, nil
	}
	return models.Board{}, ErrBoardDoesNotExist
}

// GetBoardAndUsers returns a flat structure of mock BoardAndUser domain types.
func (r *mockRepository) GetBoardAndUsers(ctx context.Context, boardID uuid.UUID) ([]BoardAndUser, error) {
	boardAndUsers := []BoardAndUser{}
	for _, boardMembership := range r.boardMemberships {
		if boardMembership.BoardID == boardID {
			board := r.boards[boardID]
			user := r.users[boardMembership.UserID]
			boardAndUser := BoardAndUser{
				Board:           board,
				BoardMembership: boardMembership,
				User:            user,
			}
			boardAndUsers = append(boardAndUsers, boardAndUser)
		}
	}
	return boardAndUsers, nil
}

// GetInvite returns a single invite.
func (r *mockRepository) GetInvite(ctx context.Context, inviteID uuid.UUID) (models.Invite, error) {
	if invite, ok := r.invites[inviteID]; ok {
		return invite, nil
	}
	return models.Invite{}, ErrInviteDoesNotExist
}

// ListOwnedBoards returns a list of mock boards belonging to a mock user.
func (r *mockRepository) ListOwnedBoards(ctx context.Context, userID uuid.UUID) ([]models.Board, error) {
	list := []models.Board{}
	for _, board := range r.boards {
		if board.UserID == userID {
			list = append(list, board)
		}
	}
	return list, nil
}

// ListOwnedBoardAndUsers returns a list of mock owned boards and associated mock members for a given mock user.
func (r *mockRepository) ListOwnedBoardAndUsers(ctx context.Context, userID uuid.UUID) ([]BoardAndUser, error) {
	ownedBoardIDs := []uuid.UUID{}
	boardAndUsers := []BoardAndUser{}
	for _, board := range r.boards {
		if board.UserID == userID {
			ownedBoardIDs = append(ownedBoardIDs, board.ID)
		}
	}
	for _, boardMembership := range r.boardMemberships {
		if boardInList(boardMembership.BoardID, ownedBoardIDs) {
			board := r.boards[boardMembership.BoardID]
			user := r.users[boardMembership.UserID]
			boardAndUser := BoardAndUser{
				Board:           board,
				BoardMembership: boardMembership,
				User:            user,
			}
			boardAndUsers = append(boardAndUsers, boardAndUser)
		}
	}
	return boardAndUsers, nil
}

// ListSharedBoardAndUsers returns a list of mock shared boards and associated mock members for a given mock user.
func (r *mockRepository) ListSharedBoardAndUsers(ctx context.Context, userID uuid.UUID) ([]BoardAndUser, error) {
	sharedBoardIDs := []uuid.UUID{}
	boardAndUsers := []BoardAndUser{}
	for _, boardMembership := range r.boardMemberships {
		if boardMembership.UserID == userID && boardMembership.Role == models.RoleMember {
			sharedBoardIDs = append(sharedBoardIDs, boardMembership.BoardID)
		}
	}
	for _, boardMembership := range r.boardMemberships {
		if boardInList(boardMembership.BoardID, sharedBoardIDs) {
			board := r.boards[boardMembership.BoardID]
			user := r.users[boardMembership.UserID]
			boardAndUser := BoardAndUser{
				Board:           board,
				BoardMembership: boardMembership,
				User:            user,
			}
			boardAndUsers = append(boardAndUsers, boardAndUser)
		}
	}
	return boardAndUsers, nil
}

// CreateMembership creates a mock membership.
func (r *mockRepository) CreateMembership(ctx context.Context, membership models.BoardMembership) error {
	r.boardMemberships[membership.ID] = membership
	return nil
}

// DeleteBoard deletes a mock board.
func (r *mockRepository) DeleteBoard(ctx context.Context, boardID uuid.UUID) error {
	delete(r.boards, boardID)
	return nil
}

// CreateInvites create mock invites.
func (r *mockRepository) CreateInvites(ctx context.Context, invites []models.Invite) error {
	for _, invite := range invites {
		r.invites[invite.ID] = invite
	}
	return nil
}

// CreateInvites create mock invites.
func (r *mockRepository) UpdateInvite(ctx context.Context, invite models.Invite) error {
	if _, ok := r.invites[invite.ID]; ok {
		r.invites[invite.ID] = invite
		return nil
	}
	return ErrInviteDoesNotExist
}

// ListInvitesByBoard returns a list of mock board invites for a given board ID and status.
func (r *mockRepository) ListInvitesByBoard(ctx context.Context, boardID uuid.UUID) ([]models.Invite, error) {
	invites := []models.Invite{}
	for _, invite := range r.invites {
		if invite.BoardID == boardID {
			invites = append(invites, invite)
		}
	}
	return invites, nil
}

// ListInvitesByReceiver returns a list of mock board invites for a given board ID and status.
func (r *mockRepository) ListInvitesByReceiver(ctx context.Context, receiverID uuid.UUID) ([]InviteBoardSender, error) {
	inviteBoardSender := []InviteBoardSender{}
	for _, invite := range r.invites {
		if invite.ReceiverID == receiverID {
			board := r.boards[invite.BoardID]
			sender := r.users[invite.SenderID]
			inviteBoardSender = append(inviteBoardSender, InviteBoardSender{invite, board, sender})
		}
	}
	return inviteBoardSender, nil
}

func boardInList(boardID uuid.UUID, list []uuid.UUID) bool {
	for _, ID := range list {
		if ID == boardID {
			return true
		}
	}
	return false
}

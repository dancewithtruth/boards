package board

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/pkg/logger"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/google/uuid"
)

var (
	ErrInvalidBoardID       = errors.New("Board ID not in UUID format")
	ErrBoardNotFound        = errors.New("Board not found.")
	defaultBoardDescription = "This is a default description for the board. Feel free to customize it and add relevant information about your board."
)

type Service interface {
	CreateBoard(ctx context.Context, input CreateBoardInput) (models.Board, error)
	GetBoard(ctx context.Context, boardID string) (models.Board, error)
	GetBoardWithMembers(ctx context.Context, boardID string) (BoardWithMembersDTO, error)
	ListOwnedBoardsWithMembers(ctx context.Context, userID string) ([]BoardWithMembersDTO, error)
	ListSharedBoardsWithMembers(ctx context.Context, userID string) ([]BoardWithMembersDTO, error)
}

type service struct {
	repo      Repository
	validator validator.Validate
}

func NewService(repo Repository, validator validator.Validate) *service {
	return &service{
		repo:      repo,
		validator: validator,
	}
}

func (s *service) CreateBoard(ctx context.Context, input CreateBoardInput) (models.Board, error) {
	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return models.Board{}, fmt.Errorf("service: failed to parse user ID input into UUID: %w", err)
	}
	// create board name if none provided
	if input.Name == nil {
		boards, err := s.repo.ListOwnedBoards(ctx, userID)
		if err != nil {
			return models.Board{}, fmt.Errorf("service: failed to get existing boards when creating board: %w", err)
		}
		numBoards := len(boards)
		boardName := fmt.Sprintf("Board #%d", numBoards+1)
		input.Name = &boardName
	}

	// use default board description if none provided
	if input.Description == nil {
		input.Description = &defaultBoardDescription
	}

	// create new board
	id := uuid.New()
	now := time.Now()
	board := models.Board{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
		UserID:      userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	err = s.repo.CreateBoard(ctx, board)
	if err != nil {
		return models.Board{}, fmt.Errorf("service: failed to create board: %w", err)
	}
	return board, nil
}

// GetBoard returns a single board for a given board ID
func (s *service) GetBoard(ctx context.Context, boardID string) (models.Board, error) {
	boardIDUUID, err := uuid.Parse(boardID)
	if err != nil {
		return models.Board{}, fmt.Errorf("service: issue parsing boardID into UUID: %w", err)
	}
	board, err := s.repo.GetBoard(ctx, boardIDUUID)
	if err != nil {
		return models.Board{}, fmt.Errorf("service: failed to get board: %w", err)
	}
	return board, nil
}

// GetBoardWithMembers returns a single board with a list of associated members
func (s *service) GetBoardWithMembers(ctx context.Context, boardID string) (BoardWithMembersDTO, error) {
	logger := logger.FromContext(ctx)
	boardIDUUID, err := uuid.Parse(boardID)
	if err != nil {
		logger.Errorf("service: failed to parse boardID")
		return BoardWithMembersDTO{}, ErrInvalidBoardID
	}
	rows, err := s.repo.GetBoardAndUsers(ctx, boardIDUUID)
	if err != nil {
		return BoardWithMembersDTO{}, fmt.Errorf("service: failed to get board with members: %w", err)
	}
	list := toBoardWithMembersDTO(rows)
	if len(list) == 0 {
		return BoardWithMembersDTO{}, ErrBoardNotFound
	}
	return list[0], nil
}

// ListOwnedBoardsWithMembers returns a list of boards that belong to a user
func (s *service) ListOwnedBoards(ctx context.Context, userID string) ([]models.Board, error) {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("service: issue parsing userID into UUID: %w", err)
	}
	boards, err := s.repo.ListOwnedBoards(ctx, userIDUUID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get boards by user ID: %w", err)
	}
	return boards, nil
}

// ListOwnedBoardsWithMembers returns a list of boards that belong to a user along with a list of board members
func (s *service) ListOwnedBoardsWithMembers(ctx context.Context, userID string) ([]BoardWithMembersDTO, error) {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("service: issue parsing userID into UUID: %w", err)
	}
	rows, err := s.repo.ListOwnedBoardAndUsers(ctx, userIDUUID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get boards by user ID: %w", err)
	}
	// transform rows into DTO
	list := toBoardWithMembersDTO(rows)
	return list, nil
}

// ListSharedBoardsWithMembers returns a list of boards that belong to a user along with a list of board members
func (s *service) ListSharedBoardsWithMembers(ctx context.Context, userID string) ([]BoardWithMembersDTO, error) {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("service: issue parsing userID into UUID: %w", err)
	}
	rows, err := s.repo.ListSharedBoardAndUsers(ctx, userIDUUID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get boards by user ID: %w", err)
	}
	// transform rows into DTO
	list := toBoardWithMembersDTO(rows)
	return list, nil
}

// toBoardWithMembersDTO transforms the BoardAndUser rows into a nested DTO struct
func toBoardWithMembersDTO(rows []BoardAndUser) []BoardWithMembersDTO {
	list := []BoardWithMembersDTO{}
	boardIDToListIndex := make(map[uuid.UUID]int)
	for _, row := range rows {
		_, exists := boardIDToListIndex[row.Board.ID]
		if !exists {
			// If board does not exist in slice, then append board into slice. Before append,
			// must convert sqlc storage type into domain type
			boardIDToListIndex[row.Board.ID] = len(list)
			newItem := BoardWithMembersDTO{
				ID:          row.Board.ID,
				Name:        row.Board.Name,
				Description: row.Board.Description,
				UserID:      row.Board.UserID,
				Members:     []MemberDTO{},
				CreatedAt:   row.Board.CreatedAt,
				UpdatedAt:   row.Board.UpdatedAt,
			}
			list = append(list, newItem)
		}
		// If user and board membership record exists, append to board members field
		if row.BoardMembership != nil && row.User != nil {
			newBoardMember := MemberDTO{
				ID:    row.User.ID,
				Name:  row.User.Name,
				Email: row.User.Email,
				Membership: MembershipDTO{
					Role:      string(row.BoardMembership.Role),
					CreatedAt: row.BoardMembership.CreatedAt,
					UpdatedAt: row.BoardMembership.UpdatedAt,
				},
				CreatedAt: row.User.CreatedAt,
				UpdatedAt: row.User.UpdatedAt,
			}
			sliceIndex := boardIDToListIndex[row.Board.ID]
			itemWithNewMember := list[sliceIndex]
			itemWithNewMember.Members = append(itemWithNewMember.Members, newBoardMember)
			list[sliceIndex] = itemWithNewMember
		}
	}
	return list
}

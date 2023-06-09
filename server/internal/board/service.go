package board

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Wave-95/boards/server/internal/models"
	"github.com/Wave-95/boards/server/pkg/logger"
	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/google/uuid"
)

var (
	ErrInvalidBoardId       = errors.New("Board ID not in UUID format")
	ErrBoardNotFound        = errors.New("Board not found.")
	defaultBoardDescription = "This is a default description for the board. Feel free to customize it and add relevant information about your board."
)

type Service interface {
	CreateBoard(ctx context.Context, input CreateBoardInput) (models.Board, error)
	GetBoard(ctx context.Context, boardId string) (models.Board, error)
	GetBoardWithMembers(ctx context.Context, boardId string) (BoardWithMembersDTO, error)
	ListOwnedBoardsWithMembers(ctx context.Context, userId string) ([]BoardWithMembersDTO, error)
	ListSharedBoardsWithMembers(ctx context.Context, userId string) ([]BoardWithMembersDTO, error)
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
	userId, err := uuid.Parse(input.UserId)
	if err != nil {
		return models.Board{}, fmt.Errorf("service: failed to parse user ID input into UUID: %w", err)
	}
	// create board name if none provided
	if input.Name == nil {
		boards, err := s.repo.ListOwnedBoards(ctx, userId)
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
		Id:          id,
		Name:        input.Name,
		Description: input.Description,
		UserId:      userId,
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
func (s *service) GetBoard(ctx context.Context, boardId string) (models.Board, error) {
	boardIdUUID, err := uuid.Parse(boardId)
	if err != nil {
		return models.Board{}, fmt.Errorf("service: issue parsing boardId into UUID: %w", err)
	}
	board, err := s.repo.GetBoard(ctx, boardIdUUID)
	if err != nil {
		return models.Board{}, fmt.Errorf("service: failed to get board: %w", err)
	}
	return board, nil
}

// GetBoardWithMembers returns a single board with a list of associated members
func (s *service) GetBoardWithMembers(ctx context.Context, boardId string) (BoardWithMembersDTO, error) {
	logger := logger.FromContext(ctx)
	boardIdUUID, err := uuid.Parse(boardId)
	if err != nil {
		logger.Errorf("service: failed to parse boardId")
		return BoardWithMembersDTO{}, ErrInvalidBoardId
	}
	rows, err := s.repo.GetBoardAndUsers(ctx, boardIdUUID)
	if err != nil {
		return BoardWithMembersDTO{}, fmt.Errorf("service: failed to get board with members: %w", err)
	}
	list := ToBoardWithMembersDTO(rows)
	if len(list) == 0 {
		return BoardWithMembersDTO{}, ErrBoardNotFound
	}
	return list[0], nil
}

// ListOwnedBoardsWithMembers returns a list of boards that belong to a user
func (s *service) ListOwnedBoards(ctx context.Context, userId string) ([]models.Board, error) {
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("service: issue parsing userId into UUID: %w", err)
	}
	boards, err := s.repo.ListOwnedBoards(ctx, userIdUUID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get boards by user ID: %w", err)
	}
	return boards, nil
}

// ListOwnedBoardsWithMembers returns a list of boards that belong to a user along with a list of board members
func (s *service) ListOwnedBoardsWithMembers(ctx context.Context, userId string) ([]BoardWithMembersDTO, error) {
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("service: issue parsing userId into UUID: %w", err)
	}
	rows, err := s.repo.ListOwnedBoardAndUsers(ctx, userIdUUID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get boards by user ID: %w", err)
	}
	// transform rows into DTO
	list := ToBoardWithMembersDTO(rows)
	return list, nil
}

// ListSharedBoardsWithMembers returns a list of boards that belong to a user along with a list of board members
func (s *service) ListSharedBoardsWithMembers(ctx context.Context, userId string) ([]BoardWithMembersDTO, error) {
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("service: issue parsing userId into UUID: %w", err)
	}
	rows, err := s.repo.ListSharedBoardAndUsers(ctx, userIdUUID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get boards by user ID: %w", err)
	}
	// transform rows into DTO
	list := ToBoardWithMembersDTO(rows)
	return list, nil
}

// ToBoardWithMembersDTO transforms the BoardAndUser rows into a nested DTO struct
func ToBoardWithMembersDTO(rows []BoardAndUser) []BoardWithMembersDTO {
	list := []BoardWithMembersDTO{}
	boardIdToListIndex := make(map[uuid.UUID]int)
	for _, row := range rows {
		_, exists := boardIdToListIndex[row.Board.Id]
		if !exists {
			// If board does not exist in slice, then append board into slice. Before append,
			// must convert sqlc storage type into domain type
			boardIdToListIndex[row.Board.Id] = len(list)
			newItem := BoardWithMembersDTO{
				Id:          row.Board.Id,
				Name:        row.Board.Name,
				Description: row.Board.Description,
				UserId:      row.Board.UserId,
				Members:     []BoardMemberDTO{},
				CreatedAt:   row.Board.CreatedAt,
				UpdatedAt:   row.Board.UpdatedAt,
			}
			list = append(list, newItem)
		}
		// If user and board membership record exists, append to board members field
		if row.BoardMembership != nil && row.User != nil {
			newBoardMember := BoardMemberDTO{
				Id:    row.User.Id,
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
			sliceIndex := boardIdToListIndex[row.Board.Id]
			itemWithNewMember := list[sliceIndex]
			itemWithNewMember.Members = append(itemWithNewMember.Members, newBoardMember)
			list[sliceIndex] = itemWithNewMember
		}
	}
	return list
}

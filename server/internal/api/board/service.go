package board

import (
	"context"
	"fmt"
	"time"

	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/google/uuid"
)

var (
	defaultBoardDescription = "This is a default description for the board. Feel free to customize it and add relevant information about the purpose, goals, or any specific details related to the board."
)

type Service interface {
	CreateBoard(ctx context.Context, input *CreateBoardInput) (*Board, error)
	GetBoard(ctx context.Context, boardId string) (*Board, error)
	GetBoardsByUserId(ctx context.Context, userId string) (Boards, error)
}

type service struct {
	repo      Repository
	validator validator.Validate
}

func (s *service) CreateBoard(ctx context.Context, input *CreateBoardInput) (*Board, error) {
	// create board name if none provided
	if input.Name == nil {
		boards, err := s.repo.GetBoardsByUserId(ctx, input.UserId)
		if err != nil {
			return nil, fmt.Errorf("service: failed to get existing boards when creating board: %w", err)
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
	board := &Board{
		Id:          id,
		Name:        input.Name,
		Description: input.Description,
		UserId:      input.UserId,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	err := s.repo.CreateBoard(ctx, board)
	if err != nil {
		return nil, fmt.Errorf("service: failed to create board: %w", err)
	}
	return board, nil
}

func (s *service) GetBoard(ctx context.Context, boardId string) (*Board, error) {
	boardIdUUID, err := uuid.Parse(boardId)
	if err != nil {
		return nil, fmt.Errorf("service: issue parsing boardId into UUID: %w", err)
	}
	board, err := s.repo.GetBoard(ctx, boardIdUUID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get board: %w", err)
	}
	return board, nil
}

func (s *service) GetBoardsByUserId(ctx context.Context, userId string) (Boards, error) {
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("service: issue parsing userId into UUID: %w", err)
	}
	boards, err := s.repo.GetBoardsByUserId(ctx, userIdUUID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get boards by user ID: %w", err)
	}
	return boards, nil
}

func NewService(repo Repository, validator validator.Validate) Service {
	return &service{repo: repo, validator: validator}
}

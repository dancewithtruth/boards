package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/google/uuid"
)

type Service interface {
	CreateUser(ctx context.Context, input *CreateUserInput) (*User, error)
	GetUser(ctx context.Context, userId string) (*User, error)
}

type service struct {
	repo      Repository
	validator validator.Validate
}

func (s *service) CreateUser(ctx context.Context, input *CreateUserInput) (*User, error) {
	id := uuid.New()
	now := time.Now()
	user := &User{
		Id:        id,
		Name:      input.Name,
		Email:     input.Email,
		Password:  input.Password,
		IsGuest:   input.IsGuest,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("service: failed to create user: %w", err)
	}
	return user, nil
}

func (s *service) GetUser(ctx context.Context, userId string) (*User, error) {
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("service: issue parsing userId into UUID: %w", err)
	}
	user, err := s.repo.GetUser(ctx, userIdUUID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get user: %w", err)
	}
	return user, nil
}

func NewService(repo Repository, validator validator.Validate) Service {
	return &service{repo: repo, validator: validator}
}

package user

import (
	"context"
	"time"

	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/google/uuid"
)

type Service interface {
	CreateUser(ctx context.Context, input *CreateUserInput) (*User, error)
}

type service struct {
	repo      Repository
	validator validator.Validate
}

func (s *service) CreateUser(ctx context.Context, input *CreateUserInput) (*User, error) {
	// validate input
	if err := s.validator.Struct(input); err != nil {
		return nil, err
	}

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
		return nil, err
	}
	return user, nil
}

func NewService(repo Repository, validator validator.Validate) Service {
	return &service{repo: repo, validator: validator}
}

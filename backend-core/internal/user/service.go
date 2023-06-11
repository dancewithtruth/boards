package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/google/uuid"
)

type Service interface {
	CreateUser(ctx context.Context, input CreateUserInput) (models.User, error)
	GetUser(ctx context.Context, userId string) (models.User, error)
	ListUsersByFuzzyEmail(ctx context.Context, input ListUsersByFuzzyEmailInput) ([]models.User, error)
}

type service struct {
	repo      Repository
	validator validator.Validate
}

func NewService(repo Repository, validator validator.Validate) *service {
	return &service{repo: repo, validator: validator}
}

func (s *service) CreateUser(ctx context.Context, input CreateUserInput) (models.User, error) {
	id := uuid.New()
	now := time.Now()
	user := models.User{
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
		return models.User{}, fmt.Errorf("service: failed to create user: %w", err)
	}
	return user, nil
}

func (s *service) GetUser(ctx context.Context, userId string) (models.User, error) {
	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return models.User{}, fmt.Errorf("service: issue parsing userId into UUID: %w", err)
	}
	user, err := s.repo.GetUser(ctx, userIdUUID)
	if err != nil {
		return models.User{}, fmt.Errorf("service: failed to get user: %w", err)
	}
	return user, nil
}

func (s *service) ListUsersByFuzzyEmail(ctx context.Context, input ListUsersByFuzzyEmailInput) ([]models.User, error) {
	if err := input.Validate(); err != nil {
		return []models.User{}, fmt.Errorf("service: invalid email: %w", err)
	}
	users, err := s.repo.ListUsersByFuzzyEmail(ctx, input.Email)
	if err != nil {
		return []models.User{}, fmt.Errorf("service: failed to list users by fuzzy email: %w", err)
	}
	for _, user := range users {
		user.Password = nil
	}
	return users, nil
}

package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/google/uuid"
)

// Service is an interface that describes all the methods pertaining to the user service.
type Service interface {
	CreateUser(ctx context.Context, input CreateUserInput) (models.User, error)
	GetUser(ctx context.Context, userID string) (models.User, error)
	ListUsersByFuzzyEmail(ctx context.Context, email string) ([]models.User, error)
}

type service struct {
	repo      Repository
	validator validator.Validate
}

// NewService initializes a service struct with dependencies
func NewService(repo Repository, validator validator.Validate) *service {
	return &service{repo: repo, validator: validator}
}

// CreateUser creates and returns a new user
func (s *service) CreateUser(ctx context.Context, input CreateUserInput) (models.User, error) {
	if err := s.validator.Struct(input); err != nil {
		return models.User{}, err
	}
	id := uuid.New()
	now := time.Now()
	user := models.User{
		ID:        id,
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

// GetUser returns a single user for a given user ID
func (s *service) GetUser(ctx context.Context, userID string) (models.User, error) {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return models.User{}, fmt.Errorf("service: issue parsing userID into UUID: %w", err)
	}
	user, err := s.repo.GetUser(ctx, userIDUUID)
	if err != nil {
		return models.User{}, fmt.Errorf("service: failed to get user: %w", err)
	}
	return user, nil
}

// ListUsersByFuzzyEmail returns a list of the top 10 users ranked by email similarity
func (s *service) ListUsersByFuzzyEmail(ctx context.Context, email string) ([]models.User, error) {
	users, err := s.repo.ListUsersByFuzzyEmail(ctx, email)
	if err != nil {
		return []models.User{}, fmt.Errorf("service: failed to list users by fuzzy email: %w", err)
	}
	for _, user := range users {
		user.Password = nil
	}
	return users, nil
}

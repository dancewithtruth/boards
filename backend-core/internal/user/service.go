package user

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/pkg/security"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/google/uuid"
)

// Service is an interface that describes all the methods pertaining to the user service.
type Service interface {
	CreateUser(ctx context.Context, input CreateUserInput) (models.User, error)
	GetUser(ctx context.Context, userID string) (models.User, error)
	ListUsersByEmail(ctx context.Context, email string) ([]models.User, error)
}

type service struct {
	userRepo  Repository
	validator validator.Validate
}

// NewService initializes a service struct with dependencies
func NewService(repo Repository, validator validator.Validate) *service {
	return &service{userRepo: repo, validator: validator}
}

// CreateUser takes a user input and standardizes the user name, hashes the password (if provided), and stores the
// user details into the database.
func (s *service) CreateUser(ctx context.Context, input CreateUserInput) (models.User, error) {
	// Validate input
	if err := s.validator.Struct(input); err != nil {
		return models.User{}, fmt.Errorf("service: failed to validate input: %w", err)
	}

	// Prepare user input
	id := uuid.New()
	name := toNameCase(input.Name)
	now := time.Now()

	// Initialize user
	user := models.User{
		ID:        id,
		Name:      name,
		Email:     input.Email,
		IsGuest:   input.IsGuest,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Hash the password
	if input.Password != nil {
		hashedPassword, err := security.HashPassword(*input.Password)
		if err != nil {
			return models.User{}, fmt.Errorf("service: faled to hash password: %w", err)
		}
		user.Password = &hashedPassword
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return models.User{}, fmt.Errorf("service: failed to create user: %w", err)
	}

	user.Password = nil
	return user, nil
}

// GetUser returns a single user for a given user ID
func (s *service) GetUser(ctx context.Context, userID string) (models.User, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return models.User{}, fmt.Errorf("service: issue parsing userID into UUID: %w", err)
	}
	user, err := s.userRepo.GetUser(ctx, userUUID)
	if err != nil {
		return models.User{}, fmt.Errorf("service: failed to get user: %w", err)
	}
	return user, nil
}

// ListUsersByEmail returns a list of the top 10 users ranked by email similarity
func (s *service) ListUsersByEmail(ctx context.Context, email string) ([]models.User, error) {
	users, err := s.userRepo.ListUsersByEmail(ctx, email)
	if err != nil {
		return []models.User{}, fmt.Errorf("service: failed to list users by fuzzy email: %w", err)
	}
	for _, user := range users {
		user.Password = nil
	}
	return users, nil
}

// toNameCase creates a regular expression to match word boundaries and convert them to name case
func toNameCase(word string) string {
	re := regexp.MustCompile(`\b\w`)
	nameCase := re.ReplaceAllStringFunc(word, strings.ToUpper)
	return nameCase
}

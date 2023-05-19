package user

import (
	"errors"
	"time"

	"github.com/Wave-95/boards/server/internal/entity"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var ErrInvalidEmail = errors.New("Invalid email address")

// TODO: Add validation for password
type CreateUserInput struct {
	Name     string
	Email    *string `validate:"omitempty,email"`
	Password *string
	IsGuest  bool
}

func (input *CreateUserInput) Validate() error {
	v := validator.New()
	if err := v.Struct(input); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Email" {
				return ErrInvalidEmail
			}
		}
	}
	return nil
}

type Service interface {
	CreateUser(input CreateUserInput) (*User, error)
}

type User struct {
	entity.User
}
type service struct {
	repo Repository
}

func (s *service) CreateUser(input CreateUserInput) (*User, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	id := uuid.New()
	now := time.Now()
	user := entity.User{
		Id:        id,
		Name:      input.Name,
		Email:     input.Email,
		Password:  input.Password,
		IsGuest:   input.IsGuest,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return &User{user}, nil
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

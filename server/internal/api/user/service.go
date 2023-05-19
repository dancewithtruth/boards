package user

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var ErrInvalidEmail = errors.New("Invalid email address")

type CreateUserInput struct {
	Name     string
	Email    string `validate:"email"`
	Password string
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
	CreateUser(input CreateUserInput) error
}

type service struct {
	repo Repository
}

func (s *service) CreateUser(input CreateUserInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return nil
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

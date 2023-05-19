package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidEmail = errors.New("Invalid email address")

// func (input *CreateUserInput) Validate() error {
// 	v := validator.New()
// 	if err := v.Struct(input); err != nil {
// 		for _, err := range err.(validator.ValidationErrors) {
// 			if err.Field() == "Email" {
// 				return ErrInvalidEmail
// 			}
// 		}
// 	}
// 	return nil
// }

type Service interface {
	CreateUser(input *CreateUserInput) (*User, error)
}

type service struct {
	repo Repository
}

func (s *service) CreateUser(input *CreateUserInput) (*User, error) {
	// TODO: Validate input
	id := uuid.New()
	now := time.Now()
	user := User{
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
	return &user, nil
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

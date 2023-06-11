package user

import (
	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
)

// Inputs

type CreateUserInput struct {
	Name     string  `json:"name" validate:"required,min=2,max=12"`
	Email    *string `json:"email" validate:"omitempty,email,required"`
	Password *string `json:"password" validate:"omitempty,min=8"`
	IsGuest  bool    `json:"is_guest" validate:"omitempty,required"`
}

type ListUsersByFuzzyEmailInput struct {
	Email string `json:"email" validate:"omitempty,email,required"`
}

func (input ListUsersByFuzzyEmailInput) Validate() error {
	validator := validator.New()
	return validator.Struct(input)
}

// DTOs

type CreateUserDTO struct {
	User     models.User `json:"user"`
	JwtToken string      `json:"jwt_token"`
}

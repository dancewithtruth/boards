package user

import (
	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
)

// CreateUserInput defines the structure for requests to create a new user.
type CreateUserInput struct {
	Name     string  `json:"name" validate:"required,min=2,max=24"`
	Email    *string `json:"email" validate:"omitempty,email,required"`
	Password *string `json:"password" validate:"omitempty,min=8"`
	IsGuest  bool    `json:"is_guest" validate:"omitempty,required"`
}

// ListUsersByFuzzyEmailInput defines the structure for requests to list users by fuzzy email.
type ListUsersByFuzzyEmailInput struct {
	Email string `json:"email" validate:"email,required"`
}

// Validate validates the iput for listing users by fuzzy email search.
func (input ListUsersByFuzzyEmailInput) Validate() error {
	validator := validator.New()
	return validator.Struct(input)
}

// CreateUserDTO defines the structure of a successful create user response.
type CreateUserDTO struct {
	User     models.User `json:"user"`
	JwtToken string      `json:"jwt_token"`
}

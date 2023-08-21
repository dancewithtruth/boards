package user

import (
	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/google/uuid"
)

// CreateUserInput defines the structure for requests to create a new user.
type CreateUserInput struct {
	Name     string  `json:"name" validate:"required,min=2,max=24"`
	Email    *string `json:"email" validate:"omitempty,email,required"`
	Password *string `json:"password" validate:"omitempty,min=8"`
	IsGuest  bool    `json:"is_guest" validate:"omitempty,required"`
}

// CreateEmailVerificationInput defines the structure for requests to create a new email verification record.
type CreateEmailVerificationInput struct {
	Code   string `json:"code" validate:"required,min=4,max=24"`
	UserID string `json:"user_id" validate:"required"`
}

// VerifyEmailInput defines the structure for verifying a user's email.
type VerifyEmailInput struct {
	Code   string `json:"code" validate:"required,min=4,max=24"`
	UserID string `json:"user_id" validate:"omitempty,required"`
}

// ListUsersByEmailInput defines the structure for requests to list users by fuzzy email.
type ListUsersByEmailInput struct {
	Email string `json:"email" validate:"email,required"`
}

// Validate validates the iput for listing users by fuzzy email search.
func (input ListUsersByEmailInput) Validate() error {
	validator := validator.New()
	return validator.Struct(input)
}

// UpdateEmailVerificationInput represents the fields that can be updated for a verification record.
type UpdateEmailVerificationInput struct {
	UserID     uuid.UUID `json:"user_id"`
	IsVerified bool      `json:"is_verified"`
}

// UpdateUserVerificationInput represents the required fields to update a user's verification status.
type UpdateUserVerificationInput struct {
	UserID     uuid.UUID `json:"user_id"`
	IsVerified bool      `json:"is_verified"`
}

// CreateUserDTO defines the structure of a successful create user response.
type CreateUserDTO struct {
	User     models.User `json:"user"`
	JwtToken string      `json:"jwt_token"`
}

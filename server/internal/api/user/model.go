package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID
	Name      string
	Email     *string
	Password  *string
	IsGuest   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) ToDto() *CreateUserResponse {
	return &CreateUserResponse{
		Id:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		IsGuest:   u.IsGuest,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type CreateUserRequest struct {
	Name     string  `json:"name" validate:"required"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
	IsGuest  bool    `json:"is_guest" validate:"omitempty,required"`
}

func (r *CreateUserRequest) ToInput() *CreateUserInput {
	return &CreateUserInput{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
		IsGuest:  r.IsGuest,
	}
}

// TODO: Add validation for password
type CreateUserInput struct {
	Name     string  `validate:"required"`
	Email    *string `validate:"omitempty,email"`
	Password *string
	IsGuest  bool
}

type CreateUserResponse struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     *string   `json:"email"`
	IsGuest   bool      `json:"is_guest"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

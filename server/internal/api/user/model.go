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

func (u *User) ToDto() UserResponse {
	return UserResponse{
		Id:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		IsGuest:   u.IsGuest,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u User) ToDtoWithToken(jwtToken string) UserResponseWithToken {
	user := UserResponse{
		Id:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		IsGuest:   u.IsGuest,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	return UserResponseWithToken{
		User:     user,
		JwtToken: jwtToken,
	}
}

type CreateUserRequest struct {
	Name     string  `json:"name" validate:"required,min=2,max=12"`
	Email    *string `json:"email" validate:"omitempty,email,required"`
	Password *string `json:"password" validate:"omitempty,min=8"`
	IsGuest  bool    `json:"is_guest" validate:"omitempty,required"`
}

func (r CreateUserRequest) ToInput() CreateUserInput {
	return CreateUserInput{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
		IsGuest:  r.IsGuest,
	}
}

// TODO: Add validation for password
type CreateUserInput struct {
	Name     string
	Email    *string
	Password *string
	IsGuest  bool
}

type UserResponse struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     *string   `json:"email"`
	IsGuest   bool      `json:"is_guest"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type UserResponseWithToken struct {
	User     UserResponse `json:"user"`
	JwtToken string       `json:"jwt_token"`
}

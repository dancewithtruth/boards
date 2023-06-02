package board

import (
	"time"

	"github.com/google/uuid"
)

// Inputs
type CreateBoardInput struct {
	Name        *string `json:"name" validate:"omitempty,required,min=3,max=20"`
	Description *string `json:"description" validate:"omitempty,required,min=3,max=100"`
	UserId      string
}

// DTOs
type BoardWithMembersDTO struct {
	Id          uuid.UUID        `json:"id"`
	Name        *string          `json:"name"`
	Description *string          `json:"description"`
	UserId      uuid.UUID        `json:"user_id"`
	Members     []BoardMemberDTO `json:"members"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type BoardMemberDTO struct {
	Id         uuid.UUID     `json:"id"`
	Name       string        `json:"name"`
	Email      *string       `json:"email"`
	Membership MembershipDTO `json:"membership"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type MembershipDTO struct {
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"added_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

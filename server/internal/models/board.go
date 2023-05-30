package models

import (
	"time"

	"github.com/google/uuid"
)

type Board struct {
	Id          uuid.UUID `json:"id"`
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	UserId      uuid.UUID `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BoardMembershipRole string

const (
	RoleMember BoardMembershipRole = "MEMBER"
	RoleAdmin  BoardMembershipRole = "ADMIN"
)

type BoardMembership struct {
	Id        uuid.UUID           `json:"id"`
	BoardId   uuid.UUID           `json:"board_id"`
	UserId    uuid.UUID           `json:"user_id"`
	Role      BoardMembershipRole `json:"role"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type BoardInviteStatus string

const (
	BoardInviteStatusPending  BoardInviteStatus = "PENDING"
	BoardInviteStatusAccepted BoardInviteStatus = "ACCEPTED"
	BoardInviteStatusIgnored  BoardInviteStatus = "IGNORED"
)

type BoardInvite struct {
	Id        uuid.UUID         `json:"id"`
	BoardId   uuid.UUID         `json:"board_id"`
	UserId    uuid.UUID         `json:"user_id"`
	Status    BoardInviteStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

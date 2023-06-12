package models

import (
	"time"

	"github.com/google/uuid"
)

type Board struct {
	ID          uuid.UUID `json:"id"`
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	UserID      uuid.UUID `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BoardMembershipRole string

const (
	RoleMember BoardMembershipRole = "MEMBER"
	RoleAdmin  BoardMembershipRole = "ADMIN"
)

type BoardMembership struct {
	ID        uuid.UUID           `json:"id"`
	BoardID   uuid.UUID           `json:"board_id"`
	UserID    uuid.UUID           `json:"user_id"`
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
	ID        uuid.UUID         `json:"id"`
	BoardID   uuid.UUID         `json:"board_id"`
	UserID    uuid.UUID         `json:"user_id"`
	Status    BoardInviteStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

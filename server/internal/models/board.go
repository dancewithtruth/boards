package models

import (
	"time"

	"github.com/google/uuid"
)

type Board struct {
	Id          uuid.UUID   `json:"id"`
	Name        *string     `json:"name"`
	Description *string     `json:"description"`
	UserId      uuid.UUID   `json:"user_id"`
	Users       []BoardUser `json:"users"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type BoardUserRole string

const (
	BoardUserRoleMember BoardUserRole = "MEMBER"
	BoardUserRoleAdmin  BoardUserRole = "ADMIN"
)

type BoardUser struct {
	Id        uuid.UUID     `json:"id"`
	BoardId   uuid.UUID     `json:"board_id"`
	UserId    uuid.UUID     `json:"user_id"`
	User      User          `json:"user"`
	Role      BoardUserRole `json:"role"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
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

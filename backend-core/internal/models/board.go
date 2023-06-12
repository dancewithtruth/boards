package models

import (
	"time"

	"github.com/google/uuid"
)

// Board defines the domain model for a board entity.
type Board struct {
	ID          uuid.UUID `json:"id"`
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	UserID      uuid.UUID `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BoardMembershipRole is a custom string type to represent board membership roles.
type BoardMembershipRole string

const (
	// RoleMember represents a role of type MEMBER.
	RoleMember BoardMembershipRole = "MEMBER"
	// RoleAdmin represents a role of type ADMIN.
	RoleAdmin BoardMembershipRole = "ADMIN"
)

// BoardMembership defines the domain model for a board membership entity.
type BoardMembership struct {
	ID        uuid.UUID           `json:"id"`
	BoardID   uuid.UUID           `json:"board_id"`
	UserID    uuid.UUID           `json:"user_id"`
	Role      BoardMembershipRole `json:"role"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

// BoardInviteStatus is a custom string type to represent invite statuses.
type BoardInviteStatus string

const (
	// BoardInviteStatusPending represents an invite in the pending state.
	BoardInviteStatusPending BoardInviteStatus = "PENDING"
	// BoardInviteStatusAccepted represents an invite in the accepted state.
	BoardInviteStatusAccepted BoardInviteStatus = "ACCEPTED"
	// BoardInviteStatusIgnored represents an invite in the ignored state.
	BoardInviteStatusIgnored BoardInviteStatus = "IGNORED"
)

// BoardInvite defines the domain model for a board invite entity.
type BoardInvite struct {
	ID         uuid.UUID         `json:"id"`
	BoardID    uuid.UUID         `json:"board_id"`
	SenderID   uuid.UUID         `json:"sender_id"`
	ReceiverID uuid.UUID         `json:"receiver_id"`
	Status     BoardInviteStatus `json:"status"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

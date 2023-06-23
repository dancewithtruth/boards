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

// InviteStatus is a custom string type to represent invite statuses.
type InviteStatus string

const (
	// InviteStatusPending represents an invite in the pending state.
	InviteStatusPending InviteStatus = "PENDING"
	// InviteStatusAccepted represents an invite in the accepted state.
	InviteStatusAccepted InviteStatus = "ACCEPTED"
	// InviteStatusIgnored represents an invite in the ignored state.
	InviteStatusIgnored InviteStatus = "IGNORED"
	// InviteStatusCancelled represents an invite in the ignored state.
	InviteStatusCancelled InviteStatus = "CANCELLED"
)

// ValidInviteStatusFilter checks if the status filter is valid if it is non-empty.
func ValidInviteStatusFilter(status string) bool {
	switch status {
	case string(InviteStatusAccepted), string(InviteStatusIgnored), string(InviteStatusCancelled), string(InviteStatusPending):
		return true
	default:
		return false
	}
}

// Invite defines the domain model for a board invite entity.
type Invite struct {
	ID         uuid.UUID    `json:"id"`
	BoardID    uuid.UUID    `json:"board_id"`
	SenderID   uuid.UUID    `json:"sender_id"`
	ReceiverID uuid.UUID    `json:"receiver_id"`
	Status     InviteStatus `json:"status"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

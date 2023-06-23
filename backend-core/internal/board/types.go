package board

import (
	"time"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/google/uuid"
)

// CreateBoardInput defines the data structure for a create board request

type CreateBoardInput struct {
	Name        *string `json:"name" validate:"omitempty,required,min=3,max=20"`
	Description *string `json:"description" validate:"omitempty,required,min=3,max=100"`
	UserID      string
}

// CreateInvitesInput defines the data structure for a create board invites request.
type CreateInvitesInput struct {
	BoardID  string
	SenderID string
	Invites  []struct {
		ReceiverId string `json:"receiver_id"`
	} `json:"invites"`
}

// UpdateInviteInput defines the data structure for a update invite request.
type UpdateInviteInput struct {
	ID     string
	UserID string
	Status string `json:"status"`
}

// ListInvitesByBoardInput defines the input params for listing invites by board.
type ListInvitesByBoardInput struct {
	BoardID string
	UserID  string
	Status  string
}

// ListInvitesByReceiverInput defines the input params for listing invites by receiver.
type ListInvitesByReceiverInput struct {
	ReceiverID string
	Status     string
}

// BoardWithMembersDTO is a formatted response representing a board and its associated members.
type BoardWithMembersDTO struct {
	ID          uuid.UUID   `json:"id"`
	Name        *string     `json:"name"`
	Description *string     `json:"description"`
	UserID      uuid.UUID   `json:"user_id"`
	Members     []MemberDTO `json:"members"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// MembersDTO is a formatted response representing a board member's details.
type MemberDTO struct {
	ID         uuid.UUID     `json:"id"`
	Name       string        `json:"name"`
	Email      *string       `json:"email"`
	Membership MembershipDTO `json:"membership"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

// MembershipDTO is a formatted response representing a board member's details.
type MembershipDTO struct {
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"added_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// InviteWithBoardAndSenderDTO is a formatted response representing a board invite along with its associated board and sender details.
type InviteWithBoardAndSenderDTO struct {
	ID         uuid.UUID    `json:"id"`
	Board      models.Board `json:"board"`
	Sender     models.User  `json:"sender"`
	ReceiverID uuid.UUID    `json:"receiver_id"`
	Status     string       `json:"status"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

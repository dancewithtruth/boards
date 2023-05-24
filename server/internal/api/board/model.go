package board

import (
	"time"

	"github.com/google/uuid"
)

type Board struct {
	Id          uuid.UUID
	Name        *string
	Description *string
	UserId      uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateBoardInput struct {
	Name        *string
	Description *string
	UserId      uuid.UUID
}

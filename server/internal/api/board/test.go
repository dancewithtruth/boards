package board

import (
	"time"

	"github.com/google/uuid"
)

func NewTestBoard(userId uuid.UUID) *Board {
	name := "test board name"
	description := "test board description"
	board := &Board{
		Id:          uuid.New(),
		Name:        &name,
		Description: &description,
		UserId:      userId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return board
}

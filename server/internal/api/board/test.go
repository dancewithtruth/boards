package board

import (
	"time"

	"github.com/Wave-95/boards/server/internal/models"
	"github.com/google/uuid"
)

func NewTestBoard(userId uuid.UUID) models.Board {
	name := "test board name"
	description := "test board description"
	board := models.Board{
		Id:          uuid.New(),
		Name:        &name,
		Description: &description,
		UserId:      userId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return board
}

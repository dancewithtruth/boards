package test

import (
	"time"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/google/uuid"
)

func NewBoard(userID uuid.UUID) models.Board {
	name := "test board name"
	description := "test board description"
	board := models.Board{
		ID:          uuid.New(),
		Name:        &name,
		Description: &description,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return board
}

package board

import (
	"context"
	"testing"

	"github.com/Wave-95/boards/server/internal/models"
	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	validator := validator.New()
	mockBoardRepo := &mockRepository{make(map[uuid.UUID]models.Board)}
	boardService := NewService(mockBoardRepo, validator)
	assert.NotNil(t, boardService)
	userId := uuid.New().String()
	t.Run("Create board", func(t *testing.T) {
		t.Run("with a board without name or description", func(t *testing.T) {
			input := CreateBoardInput{
				UserId: userId,
			}
			board, err := boardService.CreateBoard(context.Background(), input)
			assert.NoError(t, err)
			assert.Equal(t, "Board #1", *board.Name)
			assert.Equal(t, defaultBoardDescription, *board.Description)
		})

		t.Run("with a board including name or description", func(t *testing.T) {
			customBoardName := "Custom Board Name"
			customBoardDescription := "Custom board description"
			input := CreateBoardInput{
				UserId:      userId,
				Name:        &customBoardName,
				Description: &customBoardDescription,
			}
			board, err := boardService.CreateBoard(context.Background(), input)
			assert.NoError(t, err)
			assert.Equal(t, customBoardName, *board.Name)
			assert.Equal(t, customBoardDescription, *board.Description)
		})
	})

	t.Run("Get board", func(t *testing.T) {
		board, ok := getFirstBoard(mockBoardRepo.boards)
		if !ok {
			assert.FailNow(t, "expected a board to exist but got none")
		}
		board, err := boardService.GetBoard(context.Background(), board.Id.String())
		assert.NoError(t, err)
		assert.NotNil(t, board)
	})

	t.Run("Get boards", func(t *testing.T) {
		boards, err := boardService.GetBoardsByUserId(context.Background(), userId)
		assert.NoError(t, err)
		assert.Greater(t, len(boards), 0)
	})
}

func getFirstBoard(m map[uuid.UUID]models.Board) (models.Board, bool) {
	for _, board := range m {
		return board, true
	}
	return models.Board{}, false
}

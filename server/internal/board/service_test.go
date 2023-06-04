package board

import (
	"context"
	"testing"

	"github.com/Wave-95/boards/server/internal/models"
	"github.com/Wave-95/boards/server/internal/test"
	"github.com/Wave-95/boards/server/pkg/validator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	validator := validator.New()
	testUser := test.NewUser()
	mockBoardRepo := NewMockRepository(make(map[uuid.UUID]models.Board))
	mockBoardRepo.AddUser(testUser)
	boardService := NewService(mockBoardRepo, validator)
	assert.NotNil(t, boardService)
	t.Run("Create board", func(t *testing.T) {
		t.Run("without name or description", func(t *testing.T) {
			input := CreateBoardInput{
				UserId: testUser.Id.String(),
			}
			board, err := boardService.CreateBoard(context.Background(), input)
			assert.NoError(t, err)
			assert.Equal(t, "Board #1", *board.Name)
			assert.Equal(t, defaultBoardDescription, *board.Description)
		})

		t.Run("with name or description", func(t *testing.T) {
			customBoardName := "Custom Board Name"
			customBoardDescription := "Custom board description"
			input := CreateBoardInput{
				UserId:      testUser.Id.String(),
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

	t.Run("List owned boards", func(t *testing.T) {
		boards, err := boardService.ListOwnedBoards(context.Background(), testUser.Id.String())
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

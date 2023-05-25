package board

import (
	"context"
	"testing"

	"github.com/Wave-95/boards/server/internal/api/user"
	"github.com/Wave-95/boards/server/internal/models"
	"github.com/Wave-95/boards/server/internal/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	db := test.DB(t)
	userRepo := user.NewRepository(db)
	testUser := setUpTestUser(t, userRepo)

	boardRepo := NewRepository(db)
	testBoard := NewTestBoard(testUser.Id)

	t.Run("Create board", func(t *testing.T) {
		err := boardRepo.CreateBoard(context.Background(), testBoard)
		assert.NoError(t, err)
	})

	t.Run("Get board", func(t *testing.T) {
		t.Run("board exists", func(t *testing.T) {
			board, err := boardRepo.GetBoard(context.Background(), testBoard.Id)
			assert.NoError(t, err)
			assert.Equal(t, testBoard.Name, board.Name)
		})

		t.Run("board does not exist", func(t *testing.T) {
			randUUID := uuid.New()
			board, err := boardRepo.GetBoard(context.Background(), randUUID)
			assert.Empty(t, board)
			assert.ErrorIs(t, err, ErrBoardDoesNotExist)
		})
	})

	t.Run("Delete board", func(t *testing.T) {
		err := boardRepo.DeleteBoard(testBoard.Id)
		assert.NoError(t, err)
	})

	t.Run("Get boards", func(t *testing.T) {
		t.Run("no boards belong to user", func(t *testing.T) {
			boards, err := boardRepo.GetBoardsByUserId(context.Background(), testUser.Id)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(boards))
		})

		t.Run("5 boards belong to user", func(t *testing.T) {
			boardsToCreate := 5
			for i := 0; i < boardsToCreate; i++ {
				err := boardRepo.CreateBoard(context.Background(), NewTestBoard(testUser.Id))
				assert.NoError(t, err, "expected to create 5 test boards")
			}
			boards, err := boardRepo.GetBoardsByUserId(context.Background(), testUser.Id)
			assert.NoError(t, err)
			assert.Equal(t, boardsToCreate, len(boards))
			for _, board := range boards {
				boardRepo.DeleteBoard(board.Id)
			}
		})
	})

	cleanUpTestUser(t, userRepo, testUser.Id)
}

func setUpTestUser(t *testing.T, userRepo user.Repository) models.User {
	testUser := user.NewTestUser()
	err := userRepo.CreateUser(context.Background(), testUser)
	if err != nil {
		assert.FailNow(t, "Could not set up test user for board testing")
	}
	return testUser
}

func cleanUpTestUser(t *testing.T, userRepo user.Repository, userId uuid.UUID) {
	err := userRepo.DeleteUser(userId)
	if err != nil {
		assert.FailNow(t, "Could not clean up test user", err)
	}
}

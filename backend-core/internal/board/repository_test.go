package board

import (
	"context"
	"testing"
	"time"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/internal/test"
	"github.com/Wave-95/boards/backend-core/internal/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	db := test.DB(t)

	userRepo := user.NewRepository(db)
	user := setupUser(t, userRepo)
	boardRepo := NewRepository(db)

	t.Cleanup(func() {
		cleanUpTestUser(t, userRepo, user.ID)
		db.Close()
	})

	t.Run("Create, get, and delete a board", func(t *testing.T) {
		board := test.NewBoard(user.ID)

		// Create
		err := boardRepo.CreateBoard(context.Background(), board)
		assert.NoError(t, err)

		// Get
		createdBoard, err := boardRepo.GetBoard(context.Background(), board.ID)
		assert.NoError(t, err)
		assert.Equal(t, board.UserID, createdBoard.UserID)

		// Delete
		err = boardRepo.DeleteBoard(context.Background(), board.ID)
		assert.NoError(t, err)
	})

	t.Run("Get board that does not exist", func(t *testing.T) {
		randUUID := uuid.New()
		board, err := boardRepo.GetBoard(context.Background(), randUUID)
		assert.Empty(t, board)
		assert.ErrorIs(t, err, errBoardDoesNotExist)
	})

	t.Run("Create a board membership and check for membership", func(t *testing.T) {
		board := test.NewBoard(user.ID)
		// Create a membership struct to be inserted
		membership := models.BoardMembership{
			ID:        uuid.New(),
			BoardID:   board.ID,
			UserID:    user.ID,
			Role:      models.RoleMember,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := boardRepo.CreateBoard(context.Background(), board)
		if err != nil {
			assert.FailNow(t, "Failed to create test board", err)
		}
		err = boardRepo.CreateMembership(context.Background(), membership)
		assert.NoError(t, err)

		// Check that member was added to board
		boardAndUsers, err := boardRepo.GetBoardAndUsers(context.Background(), board.ID)
		firstUser := boardAndUsers[0].User
		assert.Equal(t, user.ID, firstUser.ID, "Expected user to be added to board")

		// delete board
		err = boardRepo.DeleteBoard(context.Background(), board.ID)
		assert.NoError(t, err)
	})

	t.Run("List boards by user", func(t *testing.T) {
		t.Run("owned boards", func(t *testing.T) {
			board := test.NewBoard(user.ID)
			err := boardRepo.CreateBoard(context.Background(), board)
			if err != nil {
				assert.FailNow(t, "Failed to create test board", err)
			}
			boards, err := boardRepo.ListOwnedBoards(context.Background(), user.ID)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(boards))
			err = boardRepo.DeleteBoard(context.Background(), board.ID)
			assert.NoError(t, err)
		})
	})
}

func setupUser(t *testing.T, userRepo user.Repository) models.User {
	user := test.NewUser()
	err := userRepo.CreateUser(context.Background(), user)
	if err != nil {
		assert.FailNow(t, "Failed to set up test user", err)
	}
	return user
}

func cleanUpTestUser(t *testing.T, userRepo user.Repository, userID uuid.UUID) {
	err := userRepo.DeleteUser(context.Background(), userID)
	if err != nil {
		assert.FailNow(t, "Could not clean up test user", err)
	}
}

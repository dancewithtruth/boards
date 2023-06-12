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
	testUser := setUpTestUser(t, userRepo)
	boardRepo := NewRepository(db)

	t.Cleanup(func() {
		cleanUpTestUser(t, userRepo, testUser.ID)
		db.Close()
	})

	t.Run("Create, get, and delete a board", func(t *testing.T) {
		testBoard := test.NewBoard(testUser.ID)

		// Create board
		assert.NoError(t, boardRepo.CreateBoard(context.Background(), testBoard))

		// Get board
		board, err := boardRepo.GetBoard(context.Background(), testBoard.ID)
		assert.NoError(t, err)
		assert.Equal(t, testBoard.UserID, board.UserID)

		// Delete board
		assert.NoError(t, boardRepo.DeleteBoard(context.Background(), testBoard.ID))
	})

	t.Run("Get board that does not exist", func(t *testing.T) {
		randUUID := uuid.New()
		board, err := boardRepo.GetBoard(context.Background(), randUUID)
		assert.Empty(t, board)
		assert.ErrorIs(t, err, ErrBoardDoesNotExist)
	})

	t.Run("Create a board membership and check for membership", func(t *testing.T) {
		testBoard := test.NewBoard(testUser.ID)
		// Create a membership struct to be inserted
		membership := models.BoardMembership{
			ID:        uuid.New(),
			BoardID:   testBoard.ID,
			UserID:    testUser.ID,
			Role:      models.RoleMember,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		ctx := context.Background()
		assert.NoError(t, boardRepo.CreateBoard(ctx, testBoard))
		assert.NoError(t, boardRepo.CreateMembership(ctx, membership))

		// Check that member was added to board
		boardAndUsers, err := boardRepo.GetBoardAndUsers(ctx, testBoard.ID)
		firstUser := boardAndUsers[0].User
		assert.Equal(t, testUser.ID, firstUser.ID, "expected user to be added to board")

		// delete board
		err = boardRepo.DeleteBoard(context.Background(), testBoard.ID)
		assert.NoError(t, err)
	})

	t.Run("List boards by user", func(t *testing.T) {
		t.Run("user does not have any boards", func(t *testing.T) {
			boards, err := boardRepo.ListOwnedBoardAndUsers(context.Background(), testUser.ID)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(boards))
		})

		// t.Run("user has 5 boards", func(t *testing.T) {
		// 	boardsToCreate := 5
		// 	for i := 0; i < boardsToCreate; i++ {
		// 		err := boardRepo.CreateBoard(context.Background(), NewTestBoard(testUser.Id))
		// 		assert.NoError(t, err, "expected to create 5 test boards")
		// 	}
		// 	boards, err := boardRepo.ListOwnedBoardAndUsers(context.Background(), testUser.Id)
		// 	assert.NoError(t, err)
		// 	assert.Equal(t, boardsToCreate, len(boards))
		// 	for _, board := range boards {
		// 		fmt.Println(board.Id)
		// 		assert.NoError(t, boardRepo.DeleteBoard(context.Background(), board.Id))
		// 	}
		// })
	})
}

func setUpTestUser(t *testing.T, userRepo user.Repository) models.User {
	testUser := test.NewUser()
	err := userRepo.CreateUser(context.Background(), testUser)
	if err != nil {
		assert.FailNow(t, "Could not set up test user for board testing", err)
	}
	return testUser
}

func cleanUpTestUser(t *testing.T, userRepo user.Repository, userID uuid.UUID) {
	err := userRepo.DeleteUser(context.Background(), userID)
	if err != nil {
		assert.FailNow(t, "Could not clean up test user", err)
	}
}

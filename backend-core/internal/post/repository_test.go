package post

import (
	"context"
	"testing"

	"github.com/Wave-95/boards/backend-core/db"
	"github.com/Wave-95/boards/backend-core/internal/board"
	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/internal/test"
	"github.com/Wave-95/boards/backend-core/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	db := test.DB(t)
	repo := NewRepository(db)
	assert.NotNil(t, repo)

	testUser := test.NewUser()
	testBoard := test.NewBoard(testUser.ID)
	setupUserAndBoard(t, db, testUser, testBoard)

	t.Run("Create, get, update, and delete post", func(t *testing.T) {
		// Create
		postGroup := test.NewPostGroup()
		err := repo.CreatePostGroup(context.Background(), postGroup)
		if err != nil {
			assert.FailNow(t, "Failed to create post group")
		}
		testPost := test.NewPost(testBoard.ID, testUser.ID, postGroup.ID)
		err = repo.CreatePost(context.Background(), testPost)
		assert.NoError(t, err)

		// Get
		createdPost, err := repo.GetPost(context.Background(), testPost.ID)
		assert.NoError(t, err)
		assert.Equal(t, testPost.ID, createdPost.ID)

		// Update
		updatedPost := testPost
		updatedContent := "This is updated content"
		updatedPost.Content = updatedContent
		err = repo.UpdatePost(context.Background(), updatedPost)
		assert.NoError(t, err)
		updatedPost, err = repo.GetPost(context.Background(), updatedPost.ID)
		assert.Equal(t, updatedContent, updatedPost.Content)

		// Delete
		err = repo.DeletePost(context.Background(), testPost.ID)
		assert.NoError(t, err)
		_, err = repo.GetPost(context.Background(), testPost.ID)
		assert.Error(t, err)
	})
}

func setupUserAndBoard(t *testing.T, db *db.DB, testUser models.User, testBoard models.Board) {
	userRepo := user.NewRepository(db)
	boardRepo := board.NewRepository(db)

	err := userRepo.CreateUser(context.Background(), testUser)
	if err != nil {
		t.Fatalf("Failed to create user before creating post")
	}
	err = boardRepo.CreateBoard(context.Background(), testBoard)
	if err != nil {
		t.Fatalf("Failed to create board before creating post")
	}
}

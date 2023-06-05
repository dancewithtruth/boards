package post

import (
	"context"
	"testing"

	"github.com/Wave-95/boards/server/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	mockPostRepo := NewMockRepository()
	service := NewService(mockPostRepo)
	assert.NotNil(t, service)

	t.Run("Create, get, update, and delete post", func(t *testing.T) {

		// Create
		createInput := CreatePostInput{
			UserId:  uuid.New().String(),
			BoardId: uuid.New().String(),
			Content: "This is great content right here",
			PosX:    10,
			PosY:    10,
			Color:   models.PostColorLightPink,
			ZIndex:  1,
		}
		post, err := service.CreatePost(context.Background(), createInput)
		assert.NoError(t, err)
		assert.NotEmpty(t, post.Id)

		// Update
		updatedContent := "This content has been updated"
		updatedPos := 20
		updateInput := UpdatePostInput{
			Id:      post.Id.String(),
			Content: &updatedContent,
			PosX:    &updatedPos,
			PosY:    &updatedPos,
		}
		_, err = service.UpdatePost(context.Background(), updateInput)
		assert.NoError(t, err)

		// Get
		updatedPost, err := service.GetPost(context.Background(), post.Id.String())
		assert.NoError(t, err)
		assert.Equal(t, updatedContent, updatedPost.Content)
		assert.Equal(t, updatedPos, updatedPost.PosX)

		// Delete
		err = service.DeletePost(context.Background(), post.Id.String())
		assert.NoError(t, err)

	})
}

package test

import (
	"time"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/google/uuid"
)

func NewPost(boardID uuid.UUID, userID uuid.UUID) models.Post {
	postID := uuid.New()
	now := time.Now()
	testPost := models.Post{
		ID:        postID,
		BoardID:   boardID,
		UserID:    userID,
		Content:   "This is a post!",
		PosX:      10,
		PosY:      10,
		Color:     models.PostColorLightPink,
		ZIndex:    1,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return testPost
}

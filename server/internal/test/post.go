package test

import (
	"time"

	"github.com/Wave-95/boards/server/internal/models"
	"github.com/google/uuid"
)

func NewPost(boardId uuid.UUID, userId uuid.UUID) models.Post {
	postId := uuid.New()
	now := time.Now()
	testPost := models.Post{
		Id:        postId,
		BoardId:   boardId,
		UserId:    userId,
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

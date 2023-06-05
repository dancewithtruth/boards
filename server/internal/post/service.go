package post

import (
	"context"
	"time"

	"github.com/Wave-95/boards/server/internal/models"
	"github.com/Wave-95/boards/server/pkg/logger"
	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

// CreatePost takes an input, validates it, and creates a new post
func (s *service) CreatePost(ctx context.Context, input CreatePostInput) (models.Post, error) {
	logger := logger.FromContext(ctx)
	// Validate input
	err := input.Validate()
	if err != nil {
		logger.Errorf("service: failed to validate input")
		return models.Post{}, err
	}
	// Transform input into model
	postId := uuid.New()
	userId, err := uuid.Parse(input.UserId)
	boardId, err := uuid.Parse(input.BoardId)
	if err != nil {
		logger.Errorf("service: failed to parse strings into UUIDs")
		return models.Post{}, err
	}
	now := time.Now()
	post := models.Post{
		Id:        postId,
		BoardId:   boardId,
		UserId:    userId,
		Content:   input.Content,
		PosX:      input.PosX,
		PosY:      input.PosY,
		Color:     input.Color,
		Height:    input.Height,
		ZIndex:    input.ZIndex,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = s.repo.CreatePost(ctx, post)
	if err != nil {
		logger.Errorf("service: failed to create post")
		return models.Post{}, err
	}
	return post, nil
}

// CreatePost takes an input, validates it, and creates a new post
func (s *service) GetPost(ctx context.Context, postId string) (models.Post, error) {
	logger := logger.FromContext(ctx)
	// Validate input
	postIdUUID, err := uuid.Parse(postId)
	if err != nil {
		logger.Errorf("service: failed to parse postID into UUID")
		return models.Post{}, err
	}
	//TODO: Handle error in service layer
	return s.repo.GetPost(ctx, postIdUUID)
}

func (s *service) UpdatePost(ctx context.Context, input UpdatePostInput) (models.Post, error) {
	logger := logger.FromContext(ctx)
	// Validate input
	err := input.Validate()
	if err != nil {
		logger.Errorf("service: failed to validate input")
		return models.Post{}, err
	}
	// Get post
	post, err := s.GetPost(ctx, input.Id)
	if err != nil {
		logger.Errorf("service: failed to get post for update")
		return models.Post{}, err
	}

	if input.Content != nil {
		post.Content = *input.Content
	}
	if input.PosX != nil {
		post.PosX = *input.PosX
	}
	if input.PosY != nil {
		post.PosY = *input.PosY
	}
	if input.Color != nil {
		post.Color = *input.Color
	}
	if input.Height != nil {
		post.Height = *input.Height
	}
	if input.ZIndex != nil {
		post.ZIndex = *input.ZIndex
	}
	post.UpdatedAt = time.Now()

	err = s.repo.UpdatePost(ctx, post)
	if err != nil {
		logger.Errorf("service: failed to update post")
		return models.Post{}, err
	}
	return post, nil
}

// CreatePost takes an input, validates it, and creates a new post
func (s *service) DeletePost(ctx context.Context, postId string) error {
	logger := logger.FromContext(ctx)
	postIdUUID, err := uuid.Parse(postId)
	if err != nil {
		logger.Errorf("service: failed to parse post ID into UUID")
		return err
	}
	return s.repo.DeletePost(ctx, postIdUUID)
}

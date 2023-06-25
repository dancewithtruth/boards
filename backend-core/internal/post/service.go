package post

import (
	"context"
	"time"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/pkg/logger"
	"github.com/google/uuid"
)

// Service is an interface that represents all the post service capabilities.
type Service interface {
	CreatePost(ctx context.Context, input CreatePostInput) (models.Post, error)
	GetPost(ctx context.Context, postID string) (models.Post, error)
	ListPosts(ctx context.Context, boardID string) ([]models.Post, error)
	UpdatePost(ctx context.Context, input UpdatePostInput) (models.Post, error)
	DeletePost(ctx context.Context, postID string) error
}

type service struct {
	repo Repository
}

// NewService creates a service that implements the post Service interface.
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
	postID := uuid.New()
	userID, err := uuid.Parse(input.UserID)
	boardID, err := uuid.Parse(input.BoardID)
	if err != nil {
		logger.Errorf("service: failed to parse strings into UUIDs")
		return models.Post{}, err
	}
	now := time.Now()
	post := models.Post{
		ID:        postID,
		BoardID:   boardID,
		UserID:    userID,
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

// GetPost returns a single post.
func (s *service) GetPost(ctx context.Context, postID string) (models.Post, error) {
	logger := logger.FromContext(ctx)
	// Validate input
	postUUID, err := uuid.Parse(postID)
	if err != nil {
		logger.Errorf("service: failed to parse postID into UUID")
		return models.Post{}, err
	}
	//TODO: Handle error in service layer
	return s.repo.GetPost(ctx, postUUID)
}

// ListPosts returns a list of posts for a given board ID.
func (s *service) ListPosts(ctx context.Context, boardID string) ([]models.Post, error) {
	logger := logger.FromContext(ctx)
	// Validate input
	boardUUID, err := uuid.Parse(boardID)
	if err != nil {
		logger.Errorf("service: failed to parse boardID into UUID")
		return []models.Post{}, err
	}
	//TODO: Handle error in service layer
	return s.repo.ListPosts(ctx, boardUUID)
}

// UpdatePost takes an update request and applies the updates to an exisitng post.
func (s *service) UpdatePost(ctx context.Context, input UpdatePostInput) (models.Post, error) {
	logger := logger.FromContext(ctx)
	// Validate input
	err := input.Validate()
	if err != nil {
		logger.Errorf("service: failed to validate input")
		return models.Post{}, err
	}
	// Get post
	post, err := s.GetPost(ctx, input.ID)
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

// DeletePost deletes a single post.
func (s *service) DeletePost(ctx context.Context, postID string) error {
	logger := logger.FromContext(ctx)
	postUUID, err := uuid.Parse(postID)
	if err != nil {
		logger.Errorf("service: failed to parse post ID into UUID")
		return err
	}
	return s.repo.DeletePost(ctx, postUUID)
}

package post

import (
	"context"
	"fmt"
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
	if err := input.Validate(); err != nil {
		return models.Post{}, fmt.Errorf("service: failed to validate input: %w", err)
	}
	// Prepare post
	postUUID := uuid.New()
	userUUID, err := uuid.Parse(input.UserID)
	boardUUID, err := uuid.Parse(input.BoardID)
	if err != nil {
		logger.Errorf("service: failed to parse strings into UUIDs")
		return models.Post{}, err
	}
	now := time.Now()
	var postGroupUUID uuid.UUID

	// Generate post group if post group ID not provided in input
	if input.PostGroupID == "" {
		postGroup := models.PostGroup{
			ID:        uuid.New(),
			PosX:      input.PosX,
			PosY:      input.PosY,
			ZIndex:    input.ZIndex,
			CreatedAt: now,
			UpdatedAt: now,
		}
		err := s.repo.CreatePostGroup(ctx, postGroup)
		if err != nil {
			return models.Post{}, fmt.Errorf("service: failed to auto-generate post group: %w", err)
		}
		postGroupUUID = postGroup.ID
	}

	// Assign post order to 1 if order value is not provided
	if input.PostOrder == 0 {
		input.PostOrder = 1
	}

	// Create post
	post := models.Post{
		ID:          postUUID,
		BoardID:     boardUUID,
		UserID:      userUUID,
		Content:     input.Content,
		Color:       input.Color,
		Height:      input.Height,
		CreatedAt:   now,
		UpdatedAt:   now,
		PostOrder:   input.PostOrder,
		PostGroupID: postGroupUUID,
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
	if input.Color != nil {
		post.Color = *input.Color
	}
	if input.Height != nil {
		post.Height = *input.Height
	}
	if input.PostOrder != nil {
		post.PostOrder = *input.PostOrder
	}
	if input.PostGroupID != nil {
		postGroupUUID, err := uuid.Parse(*input.PostGroupID)
		if err != nil {
			logger.Errorf("service: failed to parse post group ID")
			return models.Post{}, err
		}
		post.PostGroupID = postGroupUUID
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

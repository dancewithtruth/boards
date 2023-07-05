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
	ListPostGroups(ctx context.Context, boardID string) ([]GroupWithPostsDTO, error)
	UpdatePost(ctx context.Context, input UpdatePostInput) (models.Post, error)
	DeletePost(ctx context.Context, postID string) error
	GetPostGroup(ctx context.Context, postGroupID string) (models.PostGroup, error)
	UpdatePostGroup(ctx context.Context, input UpdatePostGroupInput) (models.PostGroup, error)
	DeletePostGroup(ctx context.Context, postGroupID string) error
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
	if err != nil {
		logger.Errorf("service: failed to parse strings into UUIDs")
		return models.Post{}, err
	}
	now := time.Now()
	var postGroupUUID uuid.UUID

	// Generate post group if post group ID not provided in input
	if input.PostGroupID == "" {
		boardUUID, err := uuid.Parse(input.BoardID)
		if err != nil {
			return models.Post{}, fmt.Errorf("service: failed to parse board ID: %w", err)
		}
		postGroup := models.PostGroup{
			ID:        uuid.New(),
			BoardID:   boardUUID,
			PosX:      input.PosX,
			PosY:      input.PosY,
			ZIndex:    input.ZIndex,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err = s.repo.CreatePostGroup(ctx, postGroup); err != nil {
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
	return s.repo.GetPost(ctx, postUUID)
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

// GetPostGroup returns a single post group.
func (s *service) GetPostGroup(ctx context.Context, postGroupID string) (models.PostGroup, error) {
	// Validate input
	postGroupUUID, err := uuid.Parse(postGroupID)
	if err != nil {
		return models.PostGroup{}, fmt.Errorf("service: failed to parse post group ID into UUID: %w", err)
	}
	return s.repo.GetPostGroup(ctx, postGroupUUID)
}

// ListPostGroups returns a list of post groups and their associated posts for a given board ID.
func (s *service) ListPostGroups(ctx context.Context, boardID string) ([]GroupWithPostsDTO, error) {
	// Validate input
	boardUUID, err := uuid.Parse(boardID)
	if err != nil {
		return []GroupWithPostsDTO{}, fmt.Errorf("service: failed to parse boardID into UUID: %w", err)
	}
	rows, err := s.repo.ListPostGroups(ctx, boardUUID)
	if err != nil {
		return []GroupWithPostsDTO{}, fmt.Errorf("service: failed to list post groups by board ID: %w", err)
	}
	return toDTOListPostGroups(rows), nil
}

// UpdatePostGroup takes an update request and applies the updates to an exisitng post group.
func (s *service) UpdatePostGroup(ctx context.Context, input UpdatePostGroupInput) (models.PostGroup, error) {
	logger := logger.FromContext(ctx)
	// Validate input
	err := input.Validate()
	if err != nil {
		logger.Errorf("service: failed to validate input")
		return models.PostGroup{}, err
	}
	// Get postGroup
	postGroup, err := s.GetPostGroup(ctx, input.ID)
	if err != nil {
		logger.Errorf("service: failed to get post group for update")
		return models.PostGroup{}, err
	}

	if input.Title != nil {
		postGroup.Title = *input.Title
	}
	if input.PosX != nil {
		postGroup.PosX = *input.PosX
	}
	if input.PosY != nil {
		postGroup.PosY = *input.PosY
	}
	if input.ZIndex != nil {
		postGroup.ZIndex = *input.ZIndex
	}
	postGroup.UpdatedAt = time.Now()

	err = s.repo.UpdatePostGroup(ctx, postGroup)
	if err != nil {
		logger.Errorf("service: failed to update post group")
		return models.PostGroup{}, err
	}
	return postGroup, nil
}

// DeletePostGroup deletes a post group for a given ID.
func (s *service) DeletePostGroup(ctx context.Context, postGroupID string) error {
	postGroupUUID, err := uuid.Parse(postGroupID)
	if err != nil {
		return fmt.Errorf("service: failed to parse ID into UUID: %w", err)
	}
	return s.repo.DeletePostGroup(ctx, postGroupUUID)
}

// toDTOListPostGroups converts the repository data structure into a nested DTO structure.
func toDTOListPostGroups(rows []GroupAndPost) []GroupWithPostsDTO {
	listDTO := []GroupWithPostsDTO{}
	parentIndex := make(map[uuid.UUID]int)
	for _, row := range rows {
		// If parent does not exist in list, add it to the list
		if _, ok := parentIndex[row.PostGroup.ID]; !ok {
			parentIndex[row.PostGroup.ID] = len(parentIndex)
			item := GroupWithPostsDTO{
				ID:        row.PostGroup.ID,
				BoardID:   row.PostGroup.BoardID,
				Title:     row.PostGroup.Title,
				PosX:      row.PostGroup.PosX,
				PosY:      row.PostGroup.PosY,
				ZIndex:    row.PostGroup.ZIndex,
				Posts:     []models.Post{},
				CreatedAt: row.PostGroup.CreatedAt,
				UpdatedAt: row.PostGroup.UpdatedAt,
			}
			listDTO = append(listDTO, item)
		}
		// Nest child into parent
		index := parentIndex[row.PostGroup.ID]
		listDTO[index].Posts = append(listDTO[index].Posts, row.Post)
	}
	return listDTO
}

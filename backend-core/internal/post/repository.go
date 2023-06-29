package post

import (
	"context"
	"errors"

	"github.com/Wave-95/boards/backend-core/db"
	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	errPostNotFound = errors.New("Post not found")
)

// Repository is an interface that represents all the database capabilities for the post repository.
type Repository interface {
	CreatePost(ctx context.Context, post models.Post) error
	CreatePostGroup(ctx context.Context, post models.PostGroup) error
	GetPost(ctx context.Context, postID uuid.UUID) (models.Post, error)
	ListPosts(ctx context.Context, postID uuid.UUID) ([]models.Post, error)
	UpdatePost(ctx context.Context, post models.Post) error
	DeletePost(ctx context.Context, postID uuid.UUID) error
}

type repository struct {
	db *db.DB
	q  *db.Queries
}

// NewRepository intializes a struct that implements the Repository interface.
func NewRepository(conn *db.DB) *repository {
	q := db.New(conn)
	return &repository{db: conn, q: q}
}

// CreatePost creates a single post.
func (r *repository) CreatePost(ctx context.Context, post models.Post) error {
	arg := db.CreatePostParams{
		ID:          pgtype.UUID{Bytes: post.ID, Valid: true},
		BoardID:     pgtype.UUID{Bytes: post.BoardID, Valid: true},
		UserID:      pgtype.UUID{Bytes: post.UserID, Valid: true},
		Content:     pgtype.Text{String: post.Content, Valid: true},
		Color:       pgtype.Text{String: post.Color, Valid: true},
		Height:      pgtype.Int4{Int32: int32(post.Height), Valid: true},
		CreatedAt:   pgtype.Timestamp{Time: post.CreatedAt, Valid: true},
		UpdatedAt:   pgtype.Timestamp{Time: post.UpdatedAt, Valid: true},
		PostOrder:   pgtype.Float8{Float64: post.PostOrder, Valid: true},
		PostGroupID: pgtype.UUID{Bytes: post.PostGroupID, Valid: true},
	}
	return r.q.CreatePost(ctx, arg)
}

// CreatePostGroup creates a single post.
func (r *repository) CreatePostGroup(ctx context.Context, postGroup models.PostGroup) error {
	arg := db.CreatePostGroupParams{
		ID:        pgtype.UUID{Bytes: postGroup.ID, Valid: true},
		Title:     pgtype.Text{String: postGroup.Title, Valid: true},
		PosX:      pgtype.Int4{Int32: int32(postGroup.PosX), Valid: true},
		PosY:      pgtype.Int4{Int32: int32(postGroup.PosY), Valid: true},
		ZIndex:    pgtype.Int4{Int32: int32(postGroup.ZIndex), Valid: true},
		CreatedAt: pgtype.Timestamp{Time: postGroup.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: postGroup.UpdatedAt, Valid: true},
	}
	return r.q.CreatePostGroup(ctx, arg)
}

// GetPost returns a single post.
func (r *repository) GetPost(ctx context.Context, postID uuid.UUID) (models.Post, error) {
	postDB, err := r.q.GetPost(ctx, pgtype.UUID{Bytes: postID, Valid: true})
	if err != nil {
		return models.Post{}, err
	}
	return toPost(postDB), nil
}

// ListPosts returns a list of posts for a given board ID.
func (r *repository) ListPosts(ctx context.Context, boardID uuid.UUID) ([]models.Post, error) {
	postsDB, err := r.q.ListPosts(ctx, pgtype.UUID{Bytes: boardID, Valid: true})
	if err != nil {
		return []models.Post{}, err
	}
	posts := []models.Post{}
	for _, postDB := range postsDB {
		posts = append(posts, toPost(postDB))
	}
	return posts, nil
}

// UpdatePost takes a post model and updates an existing post.
func (r *repository) UpdatePost(ctx context.Context, post models.Post) error {
	arg := db.UpdatePostParams{
		ID:          pgtype.UUID{Bytes: post.ID, Valid: true},
		BoardID:     pgtype.UUID{Bytes: post.BoardID, Valid: true},
		UserID:      pgtype.UUID{Bytes: post.UserID, Valid: true},
		Content:     pgtype.Text{String: post.Content, Valid: true},
		Color:       pgtype.Text{String: post.Color, Valid: true},
		Height:      pgtype.Int4{Int32: int32(post.Height), Valid: true},
		CreatedAt:   pgtype.Timestamp{Time: post.CreatedAt, Valid: true},
		UpdatedAt:   pgtype.Timestamp{Time: post.UpdatedAt, Valid: true},
		PostOrder:   pgtype.Float8{Float64: post.PostOrder, Valid: true},
		PostGroupID: pgtype.UUID{Bytes: post.PostGroupID, Valid: true},
	}
	return r.q.UpdatePost(ctx, arg)
}

// DeletePost delets a single post.
func (r *repository) DeletePost(ctx context.Context, postID uuid.UUID) error {
	return r.q.DeletePost(ctx, pgtype.UUID{Bytes: postID, Valid: true})
}

// toPost maps a db post to a domain post
func toPost(postDB db.Post) models.Post {
	return models.Post{
		ID:          postDB.ID.Bytes,
		BoardID:     postDB.BoardID.Bytes,
		UserID:      postDB.UserID.Bytes,
		Content:     postDB.Content.String,
		Color:       postDB.Color.String,
		Height:      int(postDB.Height.Int32),
		CreatedAt:   postDB.CreatedAt.Time,
		UpdatedAt:   postDB.UpdatedAt.Time,
		PostOrder:   postDB.PostOrder.Float64,
		PostGroupID: postDB.PostGroupID.Bytes,
	}
}

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
		ID:        pgtype.UUID{Bytes: post.ID, Valid: true},
		BoardID:   pgtype.UUID{Bytes: post.BoardID, Valid: true},
		UserID:    pgtype.UUID{Bytes: post.UserID, Valid: true},
		Content:   pgtype.Text{String: post.Content, Valid: true},
		PosX:      pgtype.Int4{Int32: int32(post.PosX), Valid: true},
		PosY:      pgtype.Int4{Int32: int32(post.PosY), Valid: true},
		Color:     pgtype.Text{String: post.Color, Valid: true},
		Height:    pgtype.Int4{Int32: int32(post.Height), Valid: true},
		ZIndex:    pgtype.Int4{Int32: int32(post.ZIndex), Valid: true},
		CreatedAt: pgtype.Timestamp{Time: post.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: post.UpdatedAt, Valid: true},
	}
	return r.q.CreatePost(ctx, arg)
}

// GetPost returns a single post.
func (r *repository) GetPost(ctx context.Context, postID uuid.UUID) (models.Post, error) {
	postDB, err := r.q.GetPost(ctx, pgtype.UUID{Bytes: postID, Valid: true})
	if err != nil {
		return models.Post{}, err
	}
	post := models.Post{
		ID:        postDB.ID.Bytes,
		BoardID:   postDB.BoardID.Bytes,
		UserID:    postDB.UserID.Bytes,
		Content:   postDB.Content.String,
		PosX:      int(postDB.PosX.Int32),
		PosY:      int(postDB.PosY.Int32),
		Color:     postDB.Color.String,
		Height:    int(postDB.Height.Int32),
		ZIndex:    int(postDB.ZIndex.Int32),
		CreatedAt: postDB.CreatedAt.Time,
		UpdatedAt: postDB.UpdatedAt.Time,
	}
	return post, nil
}

// ListPosts returns a list of posts for a given board ID.
func (r *repository) ListPosts(ctx context.Context, boardID uuid.UUID) ([]models.Post, error) {
	postsDB, err := r.q.ListPosts(ctx, pgtype.UUID{Bytes: boardID, Valid: true})
	if err != nil {
		return []models.Post{}, err
	}
	posts := []models.Post{}
	for _, postDB := range postsDB {
		post := models.Post{
			ID:        postDB.ID.Bytes,
			BoardID:   postDB.BoardID.Bytes,
			UserID:    postDB.UserID.Bytes,
			Content:   postDB.Content.String,
			PosX:      int(postDB.PosX.Int32),
			PosY:      int(postDB.PosY.Int32),
			Color:     postDB.Color.String,
			Height:    int(postDB.Height.Int32),
			ZIndex:    int(postDB.ZIndex.Int32),
			CreatedAt: postDB.CreatedAt.Time,
			UpdatedAt: postDB.UpdatedAt.Time,
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// UpdatePost takes a post model and updates an existing post.
func (r *repository) UpdatePost(ctx context.Context, post models.Post) error {
	arg := db.UpdatePostParams{
		ID:        pgtype.UUID{Bytes: post.ID, Valid: true},
		BoardID:   pgtype.UUID{Bytes: post.BoardID, Valid: true},
		UserID:    pgtype.UUID{Bytes: post.UserID, Valid: true},
		Content:   pgtype.Text{String: post.Content, Valid: true},
		PosX:      pgtype.Int4{Int32: int32(post.PosX), Valid: true},
		PosY:      pgtype.Int4{Int32: int32(post.PosY), Valid: true},
		Color:     pgtype.Text{String: post.Color, Valid: true},
		Height:    pgtype.Int4{Int32: int32(post.Height), Valid: true},
		ZIndex:    pgtype.Int4{Int32: int32(post.ZIndex), Valid: true},
		CreatedAt: pgtype.Timestamp{Time: post.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: post.UpdatedAt, Valid: true},
	}
	return r.q.UpdatePost(ctx, arg)
}

// DeletePost delets a single post.
func (r *repository) DeletePost(ctx context.Context, postID uuid.UUID) error {
	return r.q.DeletePost(ctx, pgtype.UUID{Bytes: postID, Valid: true})
}

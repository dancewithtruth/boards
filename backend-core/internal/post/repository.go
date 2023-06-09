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
	ErrPostNotFound = errors.New("Post not found.")
)

type Repository interface {
	CreatePost(ctx context.Context, post models.Post) error
	GetPost(ctx context.Context, postId uuid.UUID) (models.Post, error)
	ListPosts(ctx context.Context, postId uuid.UUID) ([]models.Post, error)
	UpdatePost(ctx context.Context, post models.Post) error
	DeletePost(ctx context.Context, postId uuid.UUID) error
}

type repository struct {
	db *db.DB
	q  *db.Queries
}

func NewRepository(conn *db.DB) *repository {
	q := db.New(conn)
	return &repository{db: conn, q: q}
}

func (r *repository) CreatePost(ctx context.Context, post models.Post) error {
	arg := db.CreatePostParams{
		ID:        pgtype.UUID{Bytes: post.Id, Valid: true},
		BoardID:   pgtype.UUID{Bytes: post.BoardId, Valid: true},
		UserID:    pgtype.UUID{Bytes: post.UserId, Valid: true},
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

func (r *repository) GetPost(ctx context.Context, postId uuid.UUID) (models.Post, error) {
	postDB, err := r.q.GetPost(ctx, pgtype.UUID{Bytes: postId, Valid: true})
	if err != nil {
		return models.Post{}, err
	}
	post := models.Post{
		Id:        postDB.ID.Bytes,
		BoardId:   postDB.BoardID.Bytes,
		UserId:    postDB.UserID.Bytes,
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

func (r *repository) ListPosts(ctx context.Context, boardId uuid.UUID) ([]models.Post, error) {
	postsDB, err := r.q.ListPosts(ctx, pgtype.UUID{Bytes: boardId, Valid: true})
	if err != nil {
		return []models.Post{}, err
	}
	posts := []models.Post{}
	for _, postDB := range postsDB {
		post := models.Post{
			Id:        postDB.ID.Bytes,
			BoardId:   postDB.BoardID.Bytes,
			UserId:    postDB.UserID.Bytes,
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

func (r *repository) UpdatePost(ctx context.Context, post models.Post) error {
	arg := db.UpdatePostParams{
		ID:        pgtype.UUID{Bytes: post.Id, Valid: true},
		BoardID:   pgtype.UUID{Bytes: post.BoardId, Valid: true},
		UserID:    pgtype.UUID{Bytes: post.UserId, Valid: true},
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

func (r *repository) DeletePost(ctx context.Context, postId uuid.UUID) error {
	return r.q.DeletePost(ctx, pgtype.UUID{Bytes: postId, Valid: true})
}

package post

import (
	"context"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/google/uuid"
)

type mockRepository struct {
	posts      map[uuid.UUID]models.Post
	postGroups map[uuid.UUID]models.PostGroup
}

// NewMockRepository returns a mock post repository.
func NewMockRepository() *mockRepository {
	posts := make(map[uuid.UUID]models.Post)
	postGroups := make(map[uuid.UUID]models.PostGroup)
	return &mockRepository{posts: posts, postGroups: postGroups}
}

func (r *mockRepository) CreatePost(_ context.Context, post models.Post) error {
	r.posts[post.ID] = post
	return nil
}

func (r *mockRepository) CreatePostGroup(_ context.Context, postGroup models.PostGroup) error {
	r.postGroups[postGroup.ID] = postGroup
	return nil
}

func (r *mockRepository) ListPostGroups(_ context.Context, boardID uuid.UUID) ([]GroupAndPost, error) {
	list := []GroupAndPost{}
	for _, postGroup := range r.postGroups {
		if postGroup.BoardID == boardID {
			for _, post := range r.posts {
				if post.PostGroupID == postGroup.ID {
					item := GroupAndPost{
						PostGroup: postGroup,
						Post:      post,
					}
					list = append(list, item)
				}
			}
		}
	}
	return list, nil
}

func (r *mockRepository) GetPost(_ context.Context, postID uuid.UUID) (models.Post, error) {
	if post, ok := r.posts[postID]; ok {
		return post, nil
	}
	return models.Post{}, errPostNotFound
}

func (r *mockRepository) UpdatePost(_ context.Context, post models.Post) error {
	r.posts[post.ID] = post
	return nil
}

func (r *mockRepository) DeletePost(_ context.Context, postID uuid.UUID) error {
	delete(r.posts, postID)
	return nil
}

func (r *mockRepository) GetPostGroup(_ context.Context, postGroupID uuid.UUID) (models.PostGroup, error) {
	if postGroup, ok := r.postGroups[postGroupID]; ok {
		return postGroup, nil
	}
	return models.PostGroup{}, errPostNotFound
}

func (r *mockRepository) UpdatePostGroup(_ context.Context, postGroup models.PostGroup) error {
	r.postGroups[postGroup.ID] = postGroup
	return nil
}

func (r *mockRepository) DeletePostGroup(_ context.Context, postGroupID uuid.UUID) error {
	delete(r.postGroups, postGroupID)
	return nil
}

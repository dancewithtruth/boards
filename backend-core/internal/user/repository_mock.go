package user

import (
	"context"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/google/uuid"
)

type mockRepository struct {
	users map[uuid.UUID]models.User
}

// NewMockRepository returns a mock user repository with initialized fields
func NewMockRepository() *mockRepository {
	users := make(map[uuid.UUID]models.User)
	return &mockRepository{users}
}

func (r *mockRepository) CreateUser(ctx context.Context, user models.User) error {
	r.users[user.ID] = user
	return nil
}

func (r *mockRepository) GetUser(ctx context.Context, userID uuid.UUID) (models.User, error) {
	if user, ok := r.users[userID]; ok {
		return user, nil
	}
	return models.User{}, ErrUserNotFound
}

func (r *mockRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	for _, user := range r.users {
		if email == *user.Email {
			return user, nil
		}
	}
	return models.User{}, ErrUserNotFound
}

func (r *mockRepository) ListUsersByFuzzyEmail(ctx context.Context, email string) ([]models.User, error) {
	// TODO: Mock out
	return []models.User{}, nil
}

func (r *mockRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	delete(r.users, userID)
	return nil
}

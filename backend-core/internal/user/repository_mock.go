package user

import (
	"context"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/google/uuid"
)

type mockRepository struct {
	users map[uuid.UUID]models.User
}

func (r *mockRepository) CreateUser(ctx context.Context, user models.User) error {
	r.users[user.Id] = user
	return nil
}

func (r *mockRepository) GetUser(ctx context.Context, userId uuid.UUID) (models.User, error) {
	if user, ok := r.users[userId]; ok {
		return user, nil
	}
	return models.User{}, ErrUserDoesNotExist
}

func (r *mockRepository) GetUserByLogin(ctx context.Context, email, password string) (models.User, error) {
	for _, user := range r.users {
		if email == *user.Email && password == *user.Password {
			return user, nil
		}
	}
	return models.User{}, ErrUserDoesNotExist
}

func (r *mockRepository) DeleteUser(userId uuid.UUID) error {
	delete(r.users, userId)
	return nil
}

func NewMockRepository(users map[uuid.UUID]models.User) Repository {
	return &mockRepository{users}
}

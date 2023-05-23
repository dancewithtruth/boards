package user

import (
	"context"

	"github.com/google/uuid"
)

type mockRepository struct {
	users map[uuid.UUID]*User
}

func (r *mockRepository) CreateUser(ctx context.Context, user *User) error {
	r.users[user.Id] = user
	return nil
}

func (r *mockRepository) GetUser(ctx context.Context, userId uuid.UUID) (*User, error) {
	if user, ok := r.users[userId]; ok {
		return user, nil
	}
	return nil, ErrUserDoesNotExist
}

func (r *mockRepository) GetUserByLogin(ctx context.Context, email, password string) (*User, error) {
	for _, user := range r.users {
		if email == *user.Email && password == *user.Password {
			return user, nil
		}
	}
	return nil, ErrUserDoesNotExist
}

func (r *mockRepository) DeleteUser(userId uuid.UUID) error {
	delete(r.users, userId)
	return nil
}

func NewMockRepository(users map[uuid.UUID]*User) Repository {
	return &mockRepository{users}
}

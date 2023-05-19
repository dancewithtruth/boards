package user

import (
	"github.com/Wave-95/boards/server/internal/entity"
	"github.com/google/uuid"
)

type mockRepository struct {
	users map[uuid.UUID]entity.User
}

func (r *mockRepository) CreateUser(user entity.User) error {
	r.users[user.Id] = user
	return nil
}

func (r *mockRepository) DeleteUser(userId uuid.UUID) error {
	delete(r.users, userId)
	return nil
}

package user

import (
	"github.com/google/uuid"
)

type mockRepository struct {
	users map[uuid.UUID]User
}

func (r *mockRepository) CreateUser(user User) error {
	r.users[user.Id] = user
	return nil
}

func (r *mockRepository) DeleteUser(userId uuid.UUID) error {
	delete(r.users, userId)
	return nil
}

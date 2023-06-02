package user

import (
	"context"
	"testing"

	"github.com/Wave-95/boards/server/internal/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	db := test.DB(t)
	repo := NewRepository(db)

	testUser := test.NewUser()

	t.Run("Create user", func(t *testing.T) {
		err := repo.CreateUser(context.Background(), testUser)
		assert.NoError(t, err)
	})

	t.Run("Create user with non-unique email", func(t *testing.T) {
		testUserBadEmail := testUser
		testUserBadEmail.Id = uuid.New()
		err := repo.CreateUser(context.Background(), testUserBadEmail)
		assert.ErrorIs(t, err, ErrEmailAlreadyExists)
	})

	t.Run("Get user by email and password", func(t *testing.T) {
		t.Run("email and password exist", func(t *testing.T) {
			user, err := repo.GetUserByLogin(context.Background(), *testUser.Email, *testUser.Password)
			assert.NoError(t, err)
			assert.Equal(t, testUser.Email, user.Email)
		})

		t.Run("email and password do not exist", func(t *testing.T) {
			emailNotFound := "abc123@gmail.com"
			passwordNotFound := "password1111"
			user, err := repo.GetUserByLogin(context.Background(), emailNotFound, passwordNotFound)
			assert.Empty(t, user)
			assert.ErrorIs(t, err, ErrUserDoesNotExist)
		})

	})

	t.Run("Delete user", func(t *testing.T) {
		err := repo.DeleteUser(testUser.Id)
		assert.NoError(t, err)
	})
}

package user

import (
	"context"
	"testing"

	"github.com/Wave-95/boards/backend-core/internal/test"
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
		testUserBadEmail.ID = uuid.New()
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
		err := repo.DeleteUser(context.Background(), testUser.ID)
		assert.NoError(t, err)
	})

	t.Run("List users by fuzzy email", func(t *testing.T) {
		testEmails := []string{"Georgia@gmail.com", "George@gmail.com", "Georgina@gmail.com"}
		testIds := []uuid.UUID{}
		for _, email := range testEmails {
			testUser := test.NewUser(test.WithEmail(email))
			err := repo.CreateUser(context.Background(), testUser)
			if err != nil {
				assert.FailNow(t, "Issue creating test users for fuzzy search", err)
			}
			testIds = append(testIds, testUser.ID)
		}

		defer func() {
			for _, userID := range testIds {
				repo.DeleteUser(context.Background(), userID)
			}
		}()

		users, err := repo.ListUsersByFuzzyEmail(context.Background(), "George@gmail.com")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 3, "Expected at least 3 users to be returned")
		assert.Equal(t, testEmails[1], *users[0].Email, "Expected George@gmail.com to be first result")
		assert.Equal(t, testEmails[0], *users[1].Email, "Expected Georgia@gmail.com to be second result")
		assert.Equal(t, testEmails[2], *users[2].Email, "Expected Georgina@gmail.com to be third result")
	})
}

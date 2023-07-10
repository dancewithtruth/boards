package user

import (
	"context"
	"log"
	"testing"

	"github.com/Wave-95/boards/backend-core/internal/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	db := test.DB(t)
	userRepo := NewRepository(db)
	assert.NotNil(t, userRepo)

	t.Run("Create, get, and delete user", func(t *testing.T) {
		// Create
		user := test.NewUser()
		err := userRepo.CreateUser(context.Background(), user)
		assert.NoError(t, err)

		// Get
		newUser, err := userRepo.GetUser(context.Background(), user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, newUser.ID)

		// Delete
		err = userRepo.DeleteUser(context.Background(), newUser.ID)
		assert.NoError(t, err)
		_, err = userRepo.GetUser(context.Background(), newUser.ID)
		assert.ErrorIs(t, err, ErrUserNotFound, "Expected error to be returned when user is not found")
	})

	t.Run("Create user with non-unique email results in error", func(t *testing.T) {
		// Create first user
		user := test.NewUser()
		err := userRepo.CreateUser(context.Background(), user)
		if err != nil {
			assert.FailNow(t, "Failed to create test user", err)
		}

		// Create second user using non-unique email
		userWithBadEmail := test.NewUser()
		userWithBadEmail.Email = user.Email
		err = userRepo.CreateUser(context.Background(), userWithBadEmail)
		assert.ErrorIs(t, err, errEmailAlreadyExists, "Expected error to be returned when user created with non-unique email")
		err = userRepo.DeleteUser(context.Background(), user.ID)
		if err != nil {
			log.Printf("Failed to delete test user for clean up: %v", err)
		}
	})

	t.Run("Get user by email", func(t *testing.T) {
		user := test.NewUser()
		err := userRepo.CreateUser(context.Background(), user)
		if err != nil {
			assert.FailNow(t, "Failed to create test user", err)
		}
		t.Run("email exists", func(t *testing.T) {
			newUser, err := userRepo.GetUserByEmail(context.Background(), *user.Email)
			assert.NoError(t, err)
			assert.Equal(t, user.Email, newUser.Email)
		})

		t.Run("email does not exist", func(t *testing.T) {
			email := "doesnotexist@gmail.com"
			user, err := userRepo.GetUserByEmail(context.Background(), email)
			assert.Empty(t, user)
			assert.ErrorIs(t, err, ErrUserNotFound)
		})

		err = userRepo.DeleteUser(context.Background(), user.ID)
		if err != nil {
			log.Printf("Failed to delete test user for clean up: %v", err)
		}
	})

	t.Run("List users by email", func(t *testing.T) {
		testEmails := []string{"Georgia@gmail.com", "George@gmail.com", "Georgina@gmail.com"}
		testIds := []uuid.UUID{}
		for _, email := range testEmails {
			testUser := test.NewUser(test.WithEmail(email))
			err := userRepo.CreateUser(context.Background(), testUser)
			if err != nil {
				assert.FailNow(t, "Issue creating test users for fuzzy search", err)
			}
			testIds = append(testIds, testUser.ID)
		}

		users, err := userRepo.ListUsersByFuzzyEmail(context.Background(), "George@gmail.com")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 3, "Expected at least 3 users to be returned")
		assert.Equal(t, testEmails[1], *users[0].Email, "Expected George@gmail.com to be first result")
		assert.Equal(t, testEmails[0], *users[1].Email, "Expected Georgia@gmail.com to be second result")
		assert.Equal(t, testEmails[2], *users[2].Email, "Expected Georgina@gmail.com to be third result")
		for _, userID := range testIds {
			if err := userRepo.DeleteUser(context.Background(), userID); err != nil {
				log.Printf("Failed to delete user when cleaning up test: %v", err)
			}
		}
	})
}

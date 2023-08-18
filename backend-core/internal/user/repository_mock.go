package user

import (
	"context"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/google/uuid"
)

type mockRepository struct {
	users         map[uuid.UUID]models.User
	verifications map[uuid.UUID]models.Verification
}

// NewMockRepository returns a mock user repository with initialized fields
func NewMockRepository() *mockRepository {
	users := make(map[uuid.UUID]models.User)
	verifications := make(map[uuid.UUID]models.Verification)
	return &mockRepository{users, verifications}
}

func (r *mockRepository) CreateUser(ctx context.Context, user models.User) error {
	r.users[user.ID] = user
	return nil
}

func (r *mockRepository) CreateEmailVerification(ctx context.Context, verification models.Verification) error {
	r.verifications[verification.ID] = verification
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

func (r *mockRepository) GetEmailVerification(ctx context.Context, userID uuid.UUID) (models.Verification, error) {
	// TODO: Mock out
	return models.Verification{}, nil
}

func (r *mockRepository) UpdateUserVerification(ctx context.Context, input UpdateUserVerificationInput) error {
	return nil
}

func (r *mockRepository) UpdateEmailVerification(ctx context.Context, input UpdateEmailVerificationInput) error {
	return nil
}

func (r *mockRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	delete(r.users, userID)
	return nil
}

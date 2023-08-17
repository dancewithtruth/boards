package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Wave-95/boards/backend-core/db"
	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	errEmailAlreadyExists = errors.New("User with this email already exists")
	//ErrUserNotFound is an error that occurs when a user cannot be found.
	ErrUserNotFound = errors.New("User does not exist")
	//ErrVerificationNotFound is an error that occurs when a verification cannot be found.
	ErrVerificationNotFound = errors.New("Verification does not exist")
)

// Repository represents a set of methods to interact with the database for user related responsibilities.
type Repository interface {
	CreateUser(ctx context.Context, user models.User) error
	CreateEmailVerification(ctx context.Context, verification models.Verification) error

	GetUser(ctx context.Context, userID uuid.UUID) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	ListUsersByFuzzyEmail(ctx context.Context, email string) ([]models.User, error)
	GetEmailVerification(ctx context.Context, userID uuid.UUID) (models.Verification, error)

	UpdateUserVerification(ctx context.Context, input UpdateUserVerificationInput) error
	UpdateEmailVerification(ctx context.Context, input UpdateEmailVerificationInput) error

	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

// repository implements the Repository interface.
type repository struct {
	db *db.DB
	q  *db.Queries
}

// NewRepository returns a repository struct with database and sqlc query capabilities.
func NewRepository(conn *db.DB) *repository {
	q := db.New(conn)
	return &repository{
		db: conn,
		q:  q,
	}
}

// CreateUser inserts a new user record into the users table.
func (r *repository) CreateUser(ctx context.Context, user models.User) error {
	// Prepare user for db insert
	arg := db.CreateUserParams{
		ID:        pgtype.UUID{Bytes: user.ID, Valid: true},
		Name:      pgtype.Text{String: user.Name, Valid: true},
		IsGuest:   pgtype.Bool{Bool: user.IsGuest, Valid: true},
		CreatedAt: pgtype.Timestamp{Time: user.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: user.UpdatedAt, Valid: true},
	}

	// Assign email and/or password if provided
	if user.Email != nil {
		arg.Email = pgtype.Text{String: *user.Email, Valid: true}
	}
	if user.Password != nil {
		arg.Password = pgtype.Text{String: *user.Password, Valid: true}
	}

	// Insert the user into db
	if err := r.q.CreateUser(ctx, arg); err != nil {
		if emailAlreadyExists(err) {
			return errEmailAlreadyExists
		}
		return fmt.Errorf("repository: failed to create user: %w", err)
	}
	return nil
}

// CreateEmailVerification inserts a new verification record into the email_verifications table.
func (r *repository) CreateEmailVerification(ctx context.Context, verification models.Verification) error {
	// Prepare verification for db insert
	arg := db.CreateEmailVerificationParams{
		ID:        pgtype.UUID{Bytes: verification.ID, Valid: true},
		Code:      verification.Code,
		UserID:    pgtype.UUID{Bytes: verification.UserID, Valid: true},
		CreatedAt: pgtype.Timestamp{Time: verification.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: verification.UpdatedAt, Valid: true},
	}

	// Insert record into db
	if err := r.q.CreateEmailVerification(ctx, arg); err != nil {
		return fmt.Errorf("repository: failed to create email verification: %w", err)
	}
	return nil
}

// GetUser queries and returns a single user for a given user ID.
func (r *repository) GetUser(ctx context.Context, userID uuid.UUID) (models.User, error) {
	userDB, err := r.q.GetUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("repository: failed to get user by id: %w", err)
	}
	return toUser(userDB), nil
}

// GetUserByEmail returns a single user for a given email.
func (r *repository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	userDB, err := r.q.GetUserByEmail(ctx, pgtype.Text{String: email, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("repository: failed to get user by email: %w", err)
	}
	return toUser(userDB), nil
}

// ListUsersByFuzzyEmail uses a levenshtein query to return the top 10 matches by email.
func (r *repository) ListUsersByFuzzyEmail(ctx context.Context, email string) ([]models.User, error) {
	usersDB, err := r.q.ListUsersByFuzzyEmail(ctx, pgtype.Text{String: email, Valid: true})
	if err != nil {
		return []models.User{}, fmt.Errorf("repository: failed to list users by fuzzy email: %w", err)
	}
	users := make([]models.User, len(usersDB))
	for i, userDB := range usersDB {
		user := toUser(userDB)
		users[i] = user
	}
	return users, nil
}

func (r *repository) GetEmailVerification(ctx context.Context, userID uuid.UUID) (models.Verification, error) {
	verification, err := r.q.GetEmailVerification(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Verification{}, ErrVerificationNotFound
		}
		return models.Verification{}, fmt.Errorf("repository: failed to get user by id: %w", err)
	}
	return toVerification(verification), nil
}

// UpdateUserVerification updates a user's verification status.
func (r *repository) UpdateUserVerification(ctx context.Context, input UpdateUserVerificationInput) error {
	arg := db.UpdateUserVerificationParams{
		ID:         pgtype.UUID{Bytes: input.UserID, Valid: true},
		IsVerified: pgtype.Bool{Bool: input.IsVerified, Valid: true},
	}
	err := r.q.UpdateUserVerification(ctx, arg)
	if err != nil {
		return fmt.Errorf("repository: failed to update user's verification status: %w", err)
	}
	return nil
}

// UpdateEmailVerification updates an invite.
func (r *repository) UpdateEmailVerification(ctx context.Context, input UpdateEmailVerificationInput) error {
	arg := db.UpdateEmailVerificationParams{
		UserID:     pgtype.UUID{Bytes: input.UserID, Valid: true},
		IsVerified: pgtype.Bool{Bool: input.IsVerified, Valid: true},
	}
	err := r.q.UpdateEmailVerification(ctx, arg)
	if err != nil {
		return fmt.Errorf("repository: failed to update email verification: %w", err)
	}
	return nil
}

// DeleteUser deletes a single user for a given user ID.
func (r *repository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if err := r.q.DeleteUser(ctx, pgtype.UUID{Bytes: userID, Valid: true}); err != nil {
		return fmt.Errorf("repository: failed to delete user: %w", err)
	}
	return nil
}

// toUser is a mapper that converts a db user to a domain user.
func toUser(userDB db.User) models.User {
	return models.User{
		ID:        userDB.ID.Bytes,
		Name:      userDB.Name.String,
		Email:     &userDB.Email.String,
		Password:  &userDB.Password.String,
		IsGuest:   userDB.IsGuest.Bool,
		CreatedAt: userDB.CreatedAt.Time,
		UpdatedAt: userDB.UpdatedAt.Time,
	}
}

// toVerification is a mapper that converts a db verification to a domain verification.
func toVerification(verificationDB db.EmailVerification) models.Verification {
	return models.Verification{
		ID:         verificationDB.ID.Bytes,
		Code:       verificationDB.Code,
		UserID:     verificationDB.ID.Bytes,
		IsVerified: &verificationDB.IsVerified.Bool,
		CreatedAt:  verificationDB.CreatedAt.Time,
		UpdatedAt:  verificationDB.UpdatedAt.Time,
	}
}

// emailAlreadyExists is a helper function to check if the pg error matches a unique email constraint.
func emailAlreadyExists(err error) bool {
	pgError := &pgconn.PgError{}
	return errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation && pgError.ConstraintName == "users_email_key"
}

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
)

// Repository represents a set of methods to interact with the database for user related responsibilities.
type Repository interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUser(ctx context.Context, userID uuid.UUID) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	ListUsersByFuzzyEmail(ctx context.Context, email string) ([]models.User, error)
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
	var users []models.User
	for _, userDB := range usersDB {
		user := toUser(userDB)
		users = append(users, user)
	}
	return users, nil
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

// emailAlreadyExists is a helper function to check if the pg error matches a unique email constraint.
func emailAlreadyExists(err error) bool {
	pgError := &pgconn.PgError{}
	return errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation && pgError.ConstraintName == "users_email_key"
}

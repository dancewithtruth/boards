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
	//ErrEmailAlreadyExists is an error that occurs when a user cannot be created due to duplicate email.
	ErrEmailAlreadyExists = errors.New("User with this email already exists")
	//ErrUserNotFound is an error that occurs when a user cannot be found.
	ErrUserNotFound = errors.New("User does not exist")
)

// Repository represents a set of methods to interact with the database for user related responsibilities.
type Repository interface {
	CreateUser(ctx context.Context, user models.User) error

	GetUser(ctx context.Context, userID uuid.UUID) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)

	ListUsersByEmail(ctx context.Context, email string) ([]models.User, error)

	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type repository struct {
	db *db.DB
	q  *db.Queries
}

// NewRepository initializes a repository struct with database and query capabilities.
func NewRepository(conn *db.DB) *repository {
	q := db.New(conn)
	return &repository{
		db: conn,
		q:  q,
	}
}

// CreateUser inserts a new user record into the users table.
func (r *repository) CreateUser(ctx context.Context, user models.User) error {
	sql := "INSERT INTO users (id, name, email, password, is_guest, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.db.Exec(ctx, sql, user.ID, user.Name, user.Email, user.Password, user.IsGuest, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		pgError := &pgconn.PgError{}
		if errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation && pgError.ConstraintName == "users_email_key" {
			return ErrEmailAlreadyExists
		}
		return fmt.Errorf("repository: failed to create user: %w", err)
	}
	return nil
}

// GetUser queries and returns a single user for a given user ID.
func (r *repository) GetUser(ctx context.Context, userID uuid.UUID) (models.User, error) {
	sql := "SELECT * FROM users WHERE id = $1"
	user := models.User{}
	// TODO: make scanning more robust
	err := r.db.QueryRow(ctx, sql, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IsGuest,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("repository: failed to get user by id: %w", err)
	}
	return user, nil
}

// GetUserByEmail returns a single user for a given email.
func (r *repository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	userDB, err := r.q.GetUserByEmail(ctx, pgtype.Text{String: email, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("repository: failed to get user by login credentials: %w", err)
	}
	return toUser(userDB), nil
}

// ListUsersByEmail uses a levenshtein query to return the top 10 matches by email.
func (r *repository) ListUsersByEmail(ctx context.Context, email string) ([]models.User, error) {
	rows, err := r.q.ListUsersByEmail(ctx, pgtype.Text{String: email, Valid: true})
	if err != nil {
		return []models.User{}, fmt.Errorf("repository: failed to list users by fuzzy email: %w", err)
	}
	users := []models.User{}
	for _, row := range rows {
		row := row
		user := models.User{
			ID:        row.ID.Bytes,
			Name:      row.Name.String,
			Email:     &row.Email.String,
			Password:  &row.Password.String,
			IsGuest:   row.IsGuest.Bool,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		}
		users = append(users, user)
	}
	return users, nil
}

// DeleteUser deletes a single user for a given user ID.
func (r *repository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	sql := "DELETE from users where id = $1"
	_, err := r.db.Exec(ctx, sql, userID)
	if err != nil {
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

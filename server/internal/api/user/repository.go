package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Wave-95/boards/server/db"
	"github.com/Wave-95/boards/server/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrEmailAlreadyExists = errors.New("User with this email already exists")
	ErrUserDoesNotExist   = errors.New("User does not exist")
)

type Repository interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUser(ctx context.Context, userId uuid.UUID) (models.User, error)
	GetUserByLogin(ctx context.Context, email, password string) (models.User, error)
	DeleteUser(userId uuid.UUID) error
}

type repository struct {
	db *db.DB
	q  *db.Queries
}

func NewRepository(conn *db.DB) *repository {
	q := db.New(conn)
	return &repository{
		db: conn,
		q:  q,
	}
}

func (r *repository) CreateUser(ctx context.Context, user models.User) error {
	sql := "INSERT INTO users (id, name, email, password, is_guest, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.db.Exec(ctx, sql, user.Id, user.Name, user.Email, user.Password, user.IsGuest, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		pgError := &pgconn.PgError{}
		if errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation && pgError.ConstraintName == "users_email_key" {
			return ErrEmailAlreadyExists
		}
		return fmt.Errorf("repository: failed to create user: %w", err)
	}
	return nil
}

func (r *repository) GetUser(ctx context.Context, userId uuid.UUID) (models.User, error) {
	sql := "SELECT * FROM users WHERE id = $1"
	user := models.User{}
	// TODO: make scanning more robust
	err := r.db.QueryRow(ctx, sql, userId).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IsGuest,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, ErrUserDoesNotExist
		}
		return models.User{}, fmt.Errorf("repository: failed to get user by id: %w", err)
	}
	return user, nil
}

func (r *repository) GetUserByLogin(ctx context.Context, email, password string) (models.User, error) {
	sql := "SELECT * FROM users WHERE email = $1 and password = $2"
	user := models.User{}
	// TODO: make scanning more robust
	err := r.db.QueryRow(ctx, sql, email, password).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IsGuest,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, ErrUserDoesNotExist
		}
		return models.User{}, fmt.Errorf("repository: failed to get user by login credentials: %w", err)
	}
	return user, nil
}

func (r *repository) DeleteUser(userId uuid.UUID) error {
	ctx := context.Background()
	sql := "DELETE from users where id = $1"
	_, err := r.db.Exec(ctx, sql, userId)
	if err != nil {
		return fmt.Errorf("repository: failed to delete user: %w", err)
	}
	return nil
}

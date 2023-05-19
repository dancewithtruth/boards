package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Wave-95/boards/server/db"
	"github.com/Wave-95/boards/server/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrDatabase              = errors.New("Database error")
	ErrUniqueEmailConstraint = errors.New("User with this email already exists")
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	DeleteUser(userId uuid.UUID) error
}

type repository struct {
	db *db.DB
}

func (r *repository) CreateUser(ctx context.Context, user *User) error {
	logger := logger.FromContext(ctx)
	sql := "INSERT INTO users (id, name, email, password, is_guest, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.db.Exec(ctx, sql, user.Id, user.Name, user.Email, user.Password, user.IsGuest, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		logger.Errorf("Issue inserting new user: %v", err)
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation && e.ConstraintName == "users_email_key" {
			return ErrUniqueEmailConstraint
		}
		return fmt.Errorf("%v: %w", ErrDatabase, err)
	}
	return nil
}

func (r *repository) DeleteUser(userId uuid.UUID) error {
	ctx := context.Background()
	sql := "DELETE from users where id = $1"
	_, err := r.db.Exec(ctx, sql, userId)
	if err != nil {
		return fmt.Errorf("%v: %w", ErrDatabase, err)
	}
	return nil
}

func NewRepository(db *db.DB) Repository {
	return &repository{db}
}

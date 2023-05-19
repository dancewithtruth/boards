package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Wave-95/boards/server/db"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrDatabase              = errors.New("Database error")
	ErrUniqueEmailConstraint = errors.New("Not a unique email")
)

type Repository interface {
	CreateUser(user User) error
	DeleteUser(userId uuid.UUID) error
}

type repository struct {
	db *db.DB
}

func (r *repository) CreateUser(user User) error {
	ctx := context.Background()
	sql := "INSERT INTO users (id, name, email, password, is_guest, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.db.Exec(ctx, sql, user.Id, user.Name, user.Email, user.Password, user.IsGuest, user.CreatedAt, user.UpdatedAt)
	if err != nil {
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

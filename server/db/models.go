// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Board struct {
	ID          pgtype.UUID
	Name        pgtype.Text
	Description pgtype.Text
	UserID      pgtype.UUID
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
}

type BoardInvite struct {
	ID        pgtype.UUID
	UserID    pgtype.UUID
	BoardID   pgtype.UUID
	Status    string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type BoardMembership struct {
	ID        pgtype.UUID
	UserID    pgtype.UUID
	BoardID   pgtype.UUID
	Role      pgtype.Text
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type User struct {
	ID        pgtype.UUID
	Name      pgtype.Text
	Email     pgtype.Text
	Password  pgtype.Text
	IsGuest   pgtype.Bool
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}
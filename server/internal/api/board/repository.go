package board

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Wave-95/boards/server/db"
	"github.com/Wave-95/boards/server/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
)

var (
	ErrBoardDoesNotExist = errors.New("Board does not exist")
)

type Repository interface {
	CreateBoard(ctx context.Context, board models.Board) error
	GetBoard(ctx context.Context, boardId uuid.UUID) (models.Board, error)
	ListBoardsByUser(ctx context.Context, userId uuid.UUID) ([]models.Board, error)
	DeleteBoard(ctx context.Context, boardId uuid.UUID) error
	InsertUser(ctx context.Context, user models.BoardUser) error
}

type repository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &repository{db}
}

// CreateBoard creates a single board
func (r *repository) CreateBoard(ctx context.Context, board models.Board) error {
	sql := `
	INSERT INTO boards 
	(id, name, description, user_id, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(
		ctx,
		sql,
		board.Id,
		board.Name,
		board.Description,
		board.UserId,
		board.CreatedAt,
		board.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("repository: failed to create board: %w", err)
	}
	return nil
}

// GetBoard returns a single board which contains a (possibly nil) list of associated users
func (r *repository) GetBoard(ctx context.Context, boardId uuid.UUID) (models.Board, error) {
	// // First get the board from db and scan into board struct
	// sql := `SELECT * FROM boards WHERE id = $1`
	// var board models.Board
	// err := pgxscan.Get(ctx, r.db, &board, sql, boardId)
	// if err != nil {
	// 	if errors.Is(err, pgx.ErrNoRows) {
	// 		return models.Board{}, ErrBoardDoesNotExist
	// 	}
	// 	return models.Board{}, fmt.Errorf("repository: failed to get board by id: %w", err)
	// }

	// // Next get users associated with the board
	// var users []models.User
	// sql = `
	// SELECT u.*
	// FROM users u
	// JOIN board_members ub ON u.id = ub.user_id
	// WHERE ub.board_id = $1
	// `
	// err = pgxscan.Select(ctx, r.db, &users, sql, boardId)
	// if err != nil {
	// 	return models.Board{}, fmt.Errorf("repository: failed to get users belonging to board ID: %w", err)
	// }

	// // Finally attach users to board struct and return board
	// board.Users = users
	return models.Board{}, nil
}

// ListBoardsByUser returns a list of boards for a given user along with each board's associated users
func (r *repository) ListBoardsByUser(ctx context.Context, userId uuid.UUID) ([]models.Board, error) {
	// SQL statement joins many-to-many relation between users and boards. It returns a list of boards
	// which each can have a list of associated users
	sql := `
	SELECT 
	b.id AS "b.id", 
	b.name AS "b.name", 
	b.description AS "b.description", 
	b.user_id AS "b.user_id", 
	b.created_at AS "b.created_at", 
	b.updated_at AS "b.updated_at", 

	u.id AS "u.id",
	u.name AS "u.name",
	u.email AS "u.email",
	u.is_guest AS "u.is_guest",
	u.created_at AS "u.created_at",
	u.updated_at AS "u.updated_at",

	bu.id AS "bu.id",
	bu.role AS "bu.role",
	bu.created_at AS "bu.created_at",
	bu.updated_at AS "bu.updated_at"

	FROM boards b
	LEFT JOIN board_members bu ON b.id = bu.board_id
	LEFT JOIN users u ON bu.user_id = u.id
	WHERE b.user_id = $1
	`
	// Not every board will have associated users, so a left join could result in null user values.
	// Use a NullableUser type to guard against null values
	type NullableUser struct {
		Id        *uuid.UUID
		Name      *string
		Email     *string
		IsGuest   *bool
		CreatedAt *time.Time
		UpdatedAt *time.Time
	}
	type NullableBoardUser struct {
		Id        *uuid.UUID
		Role      *string
		CreatedAt *time.Time
		UpdatedAt *time.Time
	}
	type Row struct {
		Board     models.Board      `db:"b"`
		User      NullableUser      `db:"u"`
		BoardUser NullableBoardUser `db:"bu"`
	}
	var rows []Row
	err := pgxscan.Select(ctx, r.db, &rows, sql, userId)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to list boards by user ID: %w", err)
	}

	// There can be many rows for a single board depending on the number of associated users
	// collapse the row results into a boards slice
	var boards []models.Board
	boardIdToIndex := make(map[uuid.UUID]int) // holds index of boards slice
	for _, row := range rows {
		board := row.Board
		user := row.User
		boardUser := row.BoardUser
		boardIndex, exists := boardIdToIndex[board.Id]
		// grow the boards slice for every new board
		if !exists {
			boardIdToIndex[board.Id] = len(boards)
			boards = append(boards, board)
		}
		// If user exists, convert NullableUser into User and append to the right board
		// in the boards slice. The right board is located using the boardIdToIndex map
		if user.Id != nil {
			user := models.User{
				Id:        *row.User.Id,
				Name:      *row.User.Name,
				Email:     row.User.Email,
				IsGuest:   *row.User.IsGuest,
				CreatedAt: *row.User.CreatedAt,
				UpdatedAt: *row.User.UpdatedAt,
			}
			boardUser := models.BoardUser{
				Id:        *boardUser.Id,
				Role:      models.BoardUserRole(*boardUser.Role),
				User:      user,
				CreatedAt: *boardUser.CreatedAt,
				UpdatedAt: *boardUser.UpdatedAt,
			}
			newBoard := boards[boardIndex]
			newBoard.Users = append(board.Users, boardUser)
			boards[boardIndex] = newBoard
		}
	}
	return boards, nil
}

func (r *repository) DeleteBoard(ctx context.Context, boardId uuid.UUID) error {
	sql := `DELETE from boards WHERE id = $1`
	_, err := r.db.Exec(ctx, sql, boardId)
	if err != nil {
		return fmt.Errorf("repository: failed to delete board: %w", err)
	}
	return nil
}

func (r *repository) InsertInvites(ctx context.Context, invites []models.BoardInvite) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	sql := `
	INSERT INTO board_invites
	(id, user_id, board_id, status, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	for _, invite := range invites {
		_, err = tx.Exec(ctx, sql, invite.Id, invite.UserId, invite.BoardId, invite.Status, invite.CreatedAt, invite.UpdatedAt)
		if err != nil {
			return err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

type BoardUserRole int

const (
	UserMember BoardUserRole = iota
	UserAdmin
)

func (r BoardUserRole) String() string {
	switch r {
	case UserMember:
		return "MEMBER"
	case UserAdmin:
		return "ADMIN"
	default:
		return "MEMBER"
	}
}

func (r *repository) InsertUser(ctx context.Context, user models.BoardUser) error {
	sql := `
	INSERT INTO board_members 
	(id, user_id, board_id, role, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, sql, user.Id, user.UserId, user.BoardId, user.Role, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("repository: failed to insert user: %w", err)
	}

	return nil
}

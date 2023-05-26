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
	"github.com/jackc/pgx/v5"
)

var (
	ErrBoardDoesNotExist = errors.New("Board does not exist")
)

type Repository interface {
	CreateBoard(ctx context.Context, board models.Board) error
	GetBoard(ctx context.Context, boardId uuid.UUID) (models.Board, error)
	ListBoardsByUser(ctx context.Context, userId uuid.UUID) ([]models.Board, error)
	DeleteBoard(ctx context.Context, boardId uuid.UUID) error
	InsertUsers(ctx context.Context, boardId uuid.UUID, userIds []uuid.UUID) error
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

// GetBoard returns a single board along with a list of associated users
func (r *repository) GetBoard(ctx context.Context, boardId uuid.UUID) (models.Board, error) {
	// first get the board from db and scan into board struct
	sql := `SELECT * FROM boards WHERE id = $1`
	var board models.Board
	err := pgxscan.Get(ctx, r.db, &board, sql, boardId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Board{}, ErrBoardDoesNotExist
		}
		return models.Board{}, fmt.Errorf("repository: failed to get board by id: %w", err)
	}

	// then get users associated with the board
	var users []models.User
	sql = `
	SELECT u.*
	FROM users u
	JOIN users_boards ub ON u.id = ub.user_id
	WHERE ub.board_id = $1
	`
	err = pgxscan.Select(ctx, r.db, &users, sql, boardId)
	if err != nil {
		return models.Board{}, fmt.Errorf("repository: failed to get users belonging to board ID: %w", err)
	}

	// attach users to board struct and return board
	board.Users = users
	return board, nil
}

// ListBoardsByUser returns a list of boards for a given user along with each board's associated users
func (r *repository) ListBoardsByUser(ctx context.Context, userId uuid.UUID) ([]models.Board, error) {
	// sql statement joins many-to-many relation between users and boards. It returns a list of boards
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
	u.updated_at AS "u.updated_at"

	FROM boards b
	LEFT JOIN users_boards ub ON b.id = ub.board_id
	LEFT JOIN users u ON ub.user_id = u.id
	WHERE b.user_id = $1
	`
	// use a Row type to scan row results into
	type UserRow struct {
		Id        *uuid.UUID
		Name      *string
		Email     *string
		IsGuest   *bool
		CreatedAt *time.Time
		UpdatedAt *time.Time
	}
	type Row struct {
		Board models.Board `db:"b"`
		User  UserRow      `db:"u"`
	}
	var rows []Row
	err := pgxscan.Select(ctx, r.db, &rows, sql, userId)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to list boards by user ID: %w", err)
	}
	// combines all the rows belonging to one board into a struct and return a list of board structs
	reduce := func(rows []Row) []models.Board {
		// first assign each row into corresponding board in boardMap
		boardMap := make(map[uuid.UUID]models.Board)
		for _, row := range rows {
			//check if board already exists in map
			boardId := row.Board.Id
			board, exists := boardMap[boardId]

			//if no board in map, add to map
			if !exists {
				boardMap[boardId] = row.Board
			}

			// since the sql query is using a left join, it's possible that the board has no associated
			// users and will return null values. Here we check if user ID is null before attempting to append
			// the user to Board.User slice
			if row.User.Id != nil {
				// copy UserRow contents into models.User
				user := models.User{
					Id:        *row.User.Id,
					Name:      *row.User.Name,
					Email:     row.User.Email,
					IsGuest:   *row.User.IsGuest,
					CreatedAt: *row.User.CreatedAt,
					UpdatedAt: *row.User.UpdatedAt,
				}
				board.Users = append(board.Users, user)
				boardMap[boardId] = board
			}
		}
		// convert map into slice
		boardSlice := make([]models.Board, 0, len(boardMap))
		for _, board := range boardMap {
			boardSlice = append(boardSlice, board)
		}
		return boardSlice
	}
	return reduce(rows), nil
}

func (r *repository) DeleteBoard(ctx context.Context, boardId uuid.UUID) error {
	sql := `DELETE from boards WHERE id = $1`
	_, err := r.db.Exec(ctx, sql, boardId)
	if err != nil {
		return fmt.Errorf("repository: failed to delete board: %w", err)
	}
	return nil
}

func (r *repository) InsertUsers(ctx context.Context, boardId uuid.UUID, userIds []uuid.UUID) error {
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
	INSERT INTO users_boards 
	(id, user_id, board_id) VALUES ($1, $2, $3)
	`
	for _, userId := range userIds {
		id := uuid.New() // Generate a new UUID for each row in the users_boards table
		_, err = tx.Exec(ctx, sql, id, userId, boardId)
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

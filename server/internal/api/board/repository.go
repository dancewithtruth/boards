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
	ErrBoardsNotFound    = errors.New("No boards found")
)

type Repository interface {
	CreateBoard(ctx context.Context, board models.Board) error
	GetBoard(ctx context.Context, boardId uuid.UUID) (models.Board, error)
	GetBoardsByUserId(ctx context.Context, userId uuid.UUID) ([]models.Board, error)
	DeleteBoard(ctx context.Context, boardId uuid.UUID) error

	AddUsers(ctx context.Context, boardId uuid.UUID, userIds []uuid.UUID) error
}

type repository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateBoard(ctx context.Context, board models.Board) error {
	sql := "INSERT INTO boards (id, name, description, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := r.db.Exec(ctx, sql, board.Id, board.Name, board.Description, board.UserId, board.CreatedAt, board.UpdatedAt)
	if err != nil {
		return fmt.Errorf("repository: failed to create board: %w", err)
	}
	return nil
}

// GetBoard returns a Board object with a Users slice
func (r *repository) GetBoard(ctx context.Context, boardId uuid.UUID) (models.Board, error) {
	// get board using board ID
	sql := "SELECT * FROM boards WHERE id = $1"
	var board models.Board
	err := pgxscan.Get(ctx, r.db, &board, sql, boardId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Board{}, ErrBoardDoesNotExist
		}
		return models.Board{}, fmt.Errorf("repository: failed to get board by id: %w", err)
	}
	// get users associated with board
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
	board.Users = users
	return board, nil
}

func (r *repository) GetBoardsByUserId(ctx context.Context, userId uuid.UUID) ([]models.Board, error) {
	type User struct {
		Id        *uuid.UUID
		Name      *string
		Email     *string
		IsGuest   *bool
		CreatedAt *time.Time
		UpdatedAt *time.Time
	}
	type Result struct {
		Board models.Board `db:"b"`
		User  User         `db:"u"`
	}
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
	var results []Result
	err := pgxscan.Select(ctx, r.db, &results, sql, userId)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to get boards by user_id: %w", err)
	}

	reduce := func(results []Result) []models.Board {
		boardMap := make(map[uuid.UUID]models.Board)
		for _, result := range results {
			boardId := result.Board.Id
			//check if board already exists in map
			board, ok := boardMap[boardId]

			//if no board in map, add to map
			if !ok {
				board = models.Board{
					Id:          boardId,
					Name:        result.Board.Name,
					Description: result.Board.Description,
					UserId:      result.Board.UserId,
					CreatedAt:   result.Board.CreatedAt,
					UpdatedAt:   result.Board.UpdatedAt,
				}
				boardMap[boardId] = board
			}

			//if user is valid, add user to board
			if result.User.Id != nil {
				user := models.User{
					Id:        *result.User.Id,
					Name:      *result.User.Name,
					Email:     result.User.Email,
					IsGuest:   *result.User.IsGuest,
					CreatedAt: *result.User.CreatedAt,
					UpdatedAt: *result.User.UpdatedAt,
				}
				board.Users = append(board.Users, user)
				boardMap[boardId] = board
			}
		}
		boards := make([]models.Board, 0, len(boardMap))
		for _, board := range boardMap {
			boards = append(boards, board)
		}
		return boards
	}
	return reduce(results), nil
}

func (r *repository) DeleteBoard(ctx context.Context, boardId uuid.UUID) error {
	sql := "DELETE from boards where id = $1"
	_, err := r.db.Exec(ctx, sql, boardId)
	if err != nil {
		return fmt.Errorf("repository: failed to delete board: %w", err)
	}
	return nil
}

func (r *repository) AddUsers(ctx context.Context, boardId uuid.UUID, userIds []uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	sql := "INSERT INTO users_boards (id, user_id, board_id) VALUES ($1, $2, $3)"
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

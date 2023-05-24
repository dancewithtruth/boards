package board

import (
	"context"
	"errors"
	"fmt"

	"github.com/Wave-95/boards/server/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrBoardDoesNotExist = errors.New("Board does not exist")
	ErrBoardsNotFound    = errors.New("No boards found")
)

type Repository interface {
	CreateBoard(ctx context.Context, board Board) error
	GetBoard(ctx context.Context, boardId uuid.UUID) (Board, error)
	GetBoardsByUserId(ctx context.Context, userId uuid.UUID) ([]Board, error)
	DeleteBoard(boardId uuid.UUID) error
}

type repository struct {
	db *db.DB
}

func (r *repository) CreateBoard(ctx context.Context, board Board) error {
	sql := "INSERT INTO boards (id, name, description, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := r.db.Exec(ctx, sql, board.Id, board.Name, board.Description, board.UserId, board.CreatedAt, board.UpdatedAt)
	if err != nil {
		return fmt.Errorf("repository: failed to create board: %w", err)
	}
	return nil
}

func (r *repository) GetBoard(ctx context.Context, boardId uuid.UUID) (Board, error) {
	sql := "SELECT * FROM boards WHERE id = $1"
	board := Board{}
	// TODO: make scanning more robust
	err := r.db.QueryRow(ctx, sql, boardId).Scan(
		&board.Id,
		&board.Name,
		&board.Description,
		&board.UserId,
		&board.CreatedAt,
		&board.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Board{}, ErrBoardDoesNotExist
		}
		return Board{}, fmt.Errorf("repository: failed to get board by id: %w", err)
	}
	return board, nil
}

func (r *repository) GetBoardsByUserId(ctx context.Context, boardId uuid.UUID) ([]Board, error) {
	sql := "SELECT * FROM boards WHERE user_id = $1"
	boards := []Board{}
	// TODO: make scanning more robust
	rows, err := r.db.Query(ctx, sql, boardId)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to get boards by user_id: %w", err)
	}
	for rows.Next() {
		board := Board{}
		err := rows.Scan(
			&board.Id,
			&board.Name,
			&board.Description,
			&board.UserId,
			&board.CreatedAt,
			&board.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan rows into slice: %w", err)
		}
		boards = append(boards, board)
	}
	return boards, nil
}

func (r *repository) DeleteBoard(boardId uuid.UUID) error {
	ctx := context.Background()
	sql := "DELETE from boards where id = $1"
	_, err := r.db.Exec(ctx, sql, boardId)
	if err != nil {
		return fmt.Errorf("repository: failed to delete board: %w", err)
	}
	return nil
}

func NewRepository(db *db.DB) Repository {
	return &repository{db}
}

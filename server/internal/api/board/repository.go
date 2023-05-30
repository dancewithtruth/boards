package board

import (
	"context"
	"errors"
	"fmt"

	"github.com/Wave-95/boards/server/db"
	"github.com/Wave-95/boards/server/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrBoardDoesNotExist = errors.New("Board does not exist")
)

type Repository interface {
	CreateBoard(ctx context.Context, board models.Board) error
	GetBoard(ctx context.Context, boardId uuid.UUID) (models.Board, error)
	ListOwnedBoards(ctx context.Context, userId uuid.UUID) ([]models.Board, error)
	ListOwnedBoardAndUsers(ctx context.Context, userId uuid.UUID) ([]OwnedBoardAndUser, error)
	CreateBoardInvites(ctx context.Context, invites []models.BoardInvite) error
	CreateMembership(ctx context.Context, membership models.BoardMembership) error
	DeleteBoard(ctx context.Context, boardId uuid.UUID) error
}

type repository struct {
	db *db.DB
	q  *db.Queries
}

func NewRepository(conn *db.DB) *repository {
	q := db.New(conn)
	return &repository{conn, q}
}

// CreateBoard creates a single board
func (r *repository) CreateBoard(ctx context.Context, board models.Board) error {
	// prepare board for insert
	arg := db.CreateBoardParams{
		ID:          pgtype.UUID{Bytes: board.Id, Valid: true},
		Name:        pgtype.Text{String: *board.Name, Valid: true},
		Description: pgtype.Text{String: *board.Description, Valid: true},
		UserID:      pgtype.UUID{Bytes: board.UserId, Valid: true},
		CreatedAt:   pgtype.Timestamp{Time: board.CreatedAt, Valid: true},
		UpdatedAt:   pgtype.Timestamp{Time: board.UpdatedAt, Valid: true},
	}
	err := r.q.CreateBoard(ctx, arg)
	if err != nil {
		return fmt.Errorf("repository: failed to create board: %w", err)
	}
	return nil
}

// GetBoard returns a single board
func (r *repository) GetBoard(ctx context.Context, boardId uuid.UUID) (models.Board, error) {
	row, err := r.q.GetBoard(ctx, pgtype.UUID{Bytes: boardId, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Board{}, ErrBoardDoesNotExist
		}
		return models.Board{}, fmt.Errorf("repository: failed to get board by id: %w", err)
	}
	// convert storage type to domain type
	board := models.Board{
		Id:          row.ID.Bytes,
		Name:        &row.Name.String,
		Description: &row.Description.String,
		UserId:      row.UserID.Bytes,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
	return board, nil
}

// ListOwnedBoards returns a list of boards that belong to a user
func (r *repository) ListOwnedBoards(ctx context.Context, boardId uuid.UUID) ([]models.Board, error) {
	rows, err := r.q.ListOwnedBoards(ctx, pgtype.UUID{Bytes: boardId, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []models.Board{}, ErrBoardDoesNotExist
		}
		return []models.Board{}, fmt.Errorf("repository: failed to list boards belonging to user ID: %w", err)
	}

	// convert storage type to domain type
	list := []models.Board{}
	for _, row := range rows {
		board := models.Board{
			Id:          row.ID.Bytes,
			Name:        &row.Name.String,
			Description: &row.Description.String,
			UserId:      row.UserID.Bytes,
			CreatedAt:   row.CreatedAt.Time,
			UpdatedAt:   row.UpdatedAt.Time,
		}
		list = append(list, board)
	}
	return list, nil
}

type OwnedBoardAndUser struct {
	Board           models.Board
	BoardMembership *models.BoardMembership
	User            *models.User
}

// ListOwnedBoardAndUsers returns a list of boards that a user owns along with each board's associated users
// The SQL query uses a left join so it is possible that a board can have nullable board users
func (r *repository) ListOwnedBoardAndUsers(ctx context.Context, userId uuid.UUID) ([]OwnedBoardAndUser, error) {
	rows, err := r.q.ListOwnedBoardAndUsers(ctx, pgtype.UUID{Bytes: userId, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("repository: failed to list boards by user ID: %w", err)
	}

	// Convert storage types into domain types
	list := []OwnedBoardAndUser{}
	for _, row := range rows {
		// Board will always return
		board := models.Board{
			Id:          row.Board.ID.Bytes,
			Name:        &row.Board.Name.String,
			Description: &row.Board.Description.String,
			UserId:      row.Board.UserID.Bytes,
			CreatedAt:   row.Board.CreatedAt.Time,
			UpdatedAt:   row.Board.UpdatedAt.Time,
		}
		// Initialize item with default nil placeholders for BoardMembership and User
		item := OwnedBoardAndUser{
			Board:           board,
			BoardMembership: nil,
			User:            nil,
		}
		// Check if board membership or user exists, if so then attach to item
		if row.BoardMembership.ID.Valid {
			boardMembership := models.BoardMembership{
				Id:        row.BoardMembership.ID.Bytes,
				BoardId:   row.BoardMembership.UserID.Bytes,
				UserId:    row.BoardMembership.UserID.Bytes,
				Role:      models.BoardMembershipRole(row.BoardMembership.Role.String),
				CreatedAt: row.BoardMembership.CreatedAt.Time,
				UpdatedAt: row.BoardMembership.UpdatedAt.Time,
			}
			item.BoardMembership = &boardMembership
		}
		if row.User.ID.Valid {
			user := models.User{
				Id:        row.User.ID.Bytes,
				Name:      row.User.Name.String,
				Email:     &row.User.Name.String,
				Password:  &row.User.Name.String,
				IsGuest:   row.User.IsGuest.Bool,
				CreatedAt: row.User.CreatedAt.Time,
				UpdatedAt: row.User.UpdatedAt.Time,
			}
			item.User = &user
		}
		list = append(list, item)
	}
	return list, nil
}

// DeleteBoard delets a single board
func (r *repository) DeleteBoard(ctx context.Context, boardId uuid.UUID) error {
	err := r.q.DeleteBoard(ctx, pgtype.UUID{Bytes: boardId, Valid: true})
	if err != nil {
		return fmt.Errorf("repository: failed to delete board: %w", err)
	}
	return nil
}

// CreateBoardInvites uses a db tx to insert a list of board invites. It will rollback the tx if
// any of them fail
func (r *repository) CreateBoardInvites(ctx context.Context, invites []models.BoardInvite) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	qtx := r.q.WithTx(tx)
	for _, invite := range invites {
		// prepare invite for insert
		arg := db.CreateBoardInviteParams{
			ID:        pgtype.UUID{Bytes: invite.Id, Valid: true},
			UserID:    pgtype.UUID{Bytes: invite.UserId, Valid: true},
			BoardID:   pgtype.UUID{Bytes: invite.BoardId, Valid: true},
			Status:    string(invite.Status),
			CreatedAt: pgtype.Timestamp{Time: invite.CreatedAt, Valid: true},
			UpdatedAt: pgtype.Timestamp{Time: invite.UpdatedAt, Valid: true},
		}
		err = qtx.CreateBoardInvite(ctx, arg)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// CreateMembership creates a board membership--this is effecitvely adding a user to a board
func (r *repository) CreateMembership(ctx context.Context, membership models.BoardMembership) error {
	// prepare membership for insert
	arg := db.CreateMembershipParams{
		ID:        pgtype.UUID{Bytes: membership.Id, Valid: true},
		UserID:    pgtype.UUID{Bytes: membership.UserId, Valid: true},
		BoardID:   pgtype.UUID{Bytes: membership.BoardId, Valid: true},
		Role:      pgtype.Text{String: string(membership.Role), Valid: true},
		CreatedAt: pgtype.Timestamp{Time: membership.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: membership.UpdatedAt, Valid: true},
	}
	err := r.q.CreateMembership(ctx, arg)
	if err != nil {
		return fmt.Errorf("repository: failed to insert user: %w", err)
	}
	return nil
}

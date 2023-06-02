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
	ErrTypeNotFound      = errors.New("Type not found when transforming db to domain model")
)

type Repository interface {
	CreateBoard(ctx context.Context, board models.Board) error
	GetBoard(ctx context.Context, boardId uuid.UUID) (models.Board, error)
	GetBoardAndUsers(ctx context.Context, boardId uuid.UUID) ([]BoardAndUser, error)
	ListOwnedBoards(ctx context.Context, userId uuid.UUID) ([]models.Board, error)
	ListOwnedBoardAndUsers(ctx context.Context, userId uuid.UUID) ([]BoardAndUser, error)
	ListSharedBoardAndUsers(ctx context.Context, userId uuid.UUID) ([]BoardAndUser, error)
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

// GetBoardAndUsers returns a left join query result for a single board and its associated users
func (r *repository) GetBoardAndUsers(ctx context.Context, boardId uuid.UUID) ([]BoardAndUser, error) {
	rows, err := r.q.GetBoardAndUsers(ctx, pgtype.UUID{Bytes: boardId, Valid: true})
	if err != nil {
		return []BoardAndUser{}, fmt.Errorf("repository: failed to get board and associated users: %w", err)
	}
	// Convert storage types into domain types
	list := []BoardAndUser{}
	for _, row := range rows {
		item := BoardAndUser{
			Board:           ToBoard(row.Board),
			User:            ToUser(row.User),
			BoardMembership: ToBoardMembership(row.BoardMembership),
		}
		if err != nil {
			return nil, fmt.Errorf("repository: failed to transform db row to domain model: %w", err)
		}
		list = append(list, item)
	}
	return list, nil
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

// ListOwnedBoardAndUsers returns a list of boards that a user owns along with each board's associated members
// The SQL query uses a left join so it is possible that a board can have nullable board users
func (r *repository) ListOwnedBoardAndUsers(ctx context.Context, userId uuid.UUID) ([]BoardAndUser, error) {
	rows, err := r.q.ListOwnedBoardAndUsers(ctx, pgtype.UUID{Bytes: userId, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("repository: failed to list boards by user ID: %w", err)
	}

	// Convert db types into domain types
	list := []BoardAndUser{}
	for _, row := range rows {
		item := BoardAndUser{
			Board:           ToBoard(row.Board),
			User:            ToUser(row.User),
			BoardMembership: ToBoardMembership(row.BoardMembership),
		}
		if err != nil {
			return nil, fmt.Errorf("repository: failed to transform db row to domain model: %w", err)
		}
		list = append(list, item)
	}
	return list, nil
}

// ListSharedBoardAndUsers returns a list of boards that a user belongs to along with a list of its associated members
// The SQL query uses a left join so it is possible that a board can have nullable board users
func (r *repository) ListSharedBoardAndUsers(ctx context.Context, userId uuid.UUID) ([]BoardAndUser, error) {
	rows, err := r.q.ListSharedBoardAndUsers(ctx, pgtype.UUID{Bytes: userId, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("repository: failed to list boards by user ID: %w", err)
	}

	// Convert db types into domain types
	list := []BoardAndUser{}
	for _, row := range rows {
		item := BoardAndUser{
			Board:           ToBoard(row.Board),
			User:            ToUser(row.User),
			BoardMembership: ToBoardMembership(row.BoardMembership),
		}
		if err != nil {
			return nil, fmt.Errorf("repository: failed to transform db row to domain model: %w", err)
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
			Status:    pgtype.Text{String: string(invite.Status), Valid: true},
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

// Custom repository domain types

// BoardAndUser is a custom domain type to guard against nullable BoardMembership and
// User values in the case of left joins
type BoardAndUser struct {
	Board           *models.Board
	BoardMembership *models.BoardMembership
	User            *models.User
}

// Mappers from db model to domain model

func ToBoard(dbBoard db.Board) *models.Board {
	if dbBoard.ID.Valid {
		return &models.Board{
			Id:          dbBoard.ID.Bytes,
			Name:        &dbBoard.Name.String,
			Description: &dbBoard.Description.String,
			UserId:      dbBoard.UserID.Bytes,
			CreatedAt:   dbBoard.CreatedAt.Time,
			UpdatedAt:   dbBoard.UpdatedAt.Time,
		}
	}
	return nil
}

func ToUser(dbUser db.User) *models.User {
	if dbUser.ID.Valid {
		return &models.User{
			Id:        dbUser.ID.Bytes,
			Name:      dbUser.Name.String,
			Email:     &dbUser.Email.String,
			Password:  &dbUser.Password.String,
			IsGuest:   dbUser.IsGuest.Bool,
			CreatedAt: dbUser.CreatedAt.Time,
			UpdatedAt: dbUser.UpdatedAt.Time,
		}
	}
	return nil
}

func ToBoardMembership(dbBoardMembership db.BoardMembership) *models.BoardMembership {
	if dbBoardMembership.ID.Valid {
		return &models.BoardMembership{
			Id:        dbBoardMembership.ID.Bytes,
			BoardId:   dbBoardMembership.UserID.Bytes,
			UserId:    dbBoardMembership.UserID.Bytes,
			Role:      models.BoardMembershipRole(dbBoardMembership.Role.String),
			CreatedAt: dbBoardMembership.CreatedAt.Time,
			UpdatedAt: dbBoardMembership.UpdatedAt.Time,
		}
	}
	return nil
}

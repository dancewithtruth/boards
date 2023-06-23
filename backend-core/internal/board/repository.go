package board

import (
	"context"
	"errors"
	"fmt"

	"github.com/Wave-95/boards/backend-core/db"
	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	// ErrBoardDoesNotExist is an error that is used when a board does not exist.
	ErrBoardDoesNotExist = errors.New("Board does not exist")
	// ErrInviteDoesNotExist is an error that is used when an invite does not exist.
	ErrInviteDoesNotExist = errors.New("Invite does not exist")
	// ErrTypeNotFound is an error that is used when a type is not found during storage to domain transformation.
	ErrTypeNotFound = errors.New("Type not found when transforming db to domain model")
)

// Repository is an interface that represesnts all the capabilities for interacting with the database.
type Repository interface {
	CreateBoard(ctx context.Context, board models.Board) error
	CreateMembership(ctx context.Context, membership models.BoardMembership) error
	CreateInvites(ctx context.Context, invites []models.Invite) error

	GetBoard(ctx context.Context, boardID uuid.UUID) (models.Board, error)
	GetBoardAndUsers(ctx context.Context, boardID uuid.UUID) ([]BoardAndUser, error)
	GetInvite(ctx context.Context, inviteID uuid.UUID) (models.Invite, error)

	ListOwnedBoards(ctx context.Context, userID uuid.UUID) ([]models.Board, error)
	ListOwnedBoardAndUsers(ctx context.Context, userID uuid.UUID) ([]BoardAndUser, error)
	ListSharedBoardAndUsers(ctx context.Context, userID uuid.UUID) ([]BoardAndUser, error)
	ListInvitesByBoard(ctx context.Context, boardID uuid.UUID, status string) ([]models.Invite, error)
	ListInvitesByReceiver(ctx context.Context, receiverID uuid.UUID, status string) ([]InviteBoardSender, error)

	UpdateInvite(ctx context.Context, invite models.Invite) error

	DeleteBoard(ctx context.Context, boardID uuid.UUID) error
}

type repository struct {
	db *db.DB
	q  *db.Queries
}

// NewRepository initializes and returns a repository struct with database and query capabilities.
func NewRepository(conn *db.DB) *repository {
	q := db.New(conn)
	return &repository{conn, q}
}

// CreateBoard creates a single board.
func (r *repository) CreateBoard(ctx context.Context, board models.Board) error {
	// Prepare board for insert
	arg := db.CreateBoardParams{
		ID:          pgtype.UUID{Bytes: board.ID, Valid: true},
		Name:        pgtype.Text{String: *board.Name, Valid: true},
		Description: pgtype.Text{String: *board.Description, Valid: true},
		UserID:      pgtype.UUID{Bytes: board.UserID, Valid: true},
		CreatedAt:   pgtype.Timestamp{Time: board.CreatedAt, Valid: true},
		UpdatedAt:   pgtype.Timestamp{Time: board.UpdatedAt, Valid: true},
	}
	err := r.q.CreateBoard(ctx, arg)
	if err != nil {
		return fmt.Errorf("repository: failed to create board: %w", err)
	}
	return nil
}

// CreateMembership creates a board membership--this is effecitvely adding a user to a board.
func (r *repository) CreateMembership(ctx context.Context, membership models.BoardMembership) error {
	// prepare membership for insert
	arg := db.CreateMembershipParams{
		ID:        pgtype.UUID{Bytes: membership.ID, Valid: true},
		UserID:    pgtype.UUID{Bytes: membership.UserID, Valid: true},
		BoardID:   pgtype.UUID{Bytes: membership.BoardID, Valid: true},
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

// CreateInvites uses a db tx to insert a list of board invites. It will rollback the tx if
// any of them fail.
func (r *repository) CreateInvites(ctx context.Context, invites []models.Invite) error {
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
		arg := db.CreateInviteParams(toInviteDB(invite))
		err = qtx.CreateInvite(ctx, arg)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// GetBoard returns a single board for a given board ID.
func (r *repository) GetBoard(ctx context.Context, boardID uuid.UUID) (models.Board, error) {
	row, err := r.q.GetBoard(ctx, pgtype.UUID{Bytes: boardID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Board{}, ErrBoardDoesNotExist
		}
		return models.Board{}, fmt.Errorf("repository: failed to get board by id: %w", err)
	}
	// Convert storage type to domain type.
	board := models.Board{
		ID:          row.ID.Bytes,
		Name:        &row.Name.String,
		Description: &row.Description.String,
		UserID:      row.UserID.Bytes,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
	return board, nil
}

// GetBoardAndUsers returns a flat structure list of board and users. A BoardAndUser encapsulates Board,
// User, and Membership domain models.
func (r *repository) GetBoardAndUsers(ctx context.Context, boardID uuid.UUID) ([]BoardAndUser, error) {
	rows, err := r.q.GetBoardAndUsers(ctx, pgtype.UUID{Bytes: boardID, Valid: true})
	if err != nil {
		return []BoardAndUser{}, fmt.Errorf("repository: failed to get board and associated users: %w", err)
	}
	// Convert storage types into domain types.
	list := []BoardAndUser{}
	for _, row := range rows {
		item := BoardAndUser{
			Board:           toBoard(row.Board),
			User:            toUser(row.User),
			BoardMembership: toBoardMembership(row.BoardMembership),
		}
		if err != nil {
			return nil, fmt.Errorf("repository: failed to transform db row to domain model: %w", err)
		}
		list = append(list, item)
	}
	return list, nil
}

// GetInvite returns a single invite for a given invite ID.
func (r *repository) GetInvite(ctx context.Context, inviteID uuid.UUID) (models.Invite, error) {
	row, err := r.q.GetInvite(ctx, pgtype.UUID{Bytes: inviteID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Invite{}, ErrInviteDoesNotExist
		}
		return models.Invite{}, fmt.Errorf("repository: failed to get invite by id: %w", err)
	}
	// Convert storage type to domain type.
	invite := models.Invite{
		ID:         row.ID.Bytes,
		BoardID:    row.BoardID.Bytes,
		SenderID:   row.SenderID.Bytes,
		ReceiverID: row.ReceiverID.Bytes,
		Status:     models.InviteStatus(row.Status.String),
		CreatedAt:  row.CreatedAt.Time,
		UpdatedAt:  row.UpdatedAt.Time,
	}
	return invite, nil
}

// ListOwnedBoards returns a list of boards that belong to a user.
func (r *repository) ListOwnedBoards(ctx context.Context, boardID uuid.UUID) ([]models.Board, error) {
	rows, err := r.q.ListOwnedBoards(ctx, pgtype.UUID{Bytes: boardID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []models.Board{}, ErrBoardDoesNotExist
		}
		return []models.Board{}, fmt.Errorf("repository: failed to list boards belonging to user ID: %w", err)
	}

	// convert storage type to domain type.
	list := []models.Board{}
	for _, row := range rows {
		board := models.Board{
			ID:          row.ID.Bytes,
			Name:        &row.Name.String,
			Description: &row.Description.String,
			UserID:      row.UserID.Bytes,
			CreatedAt:   row.CreatedAt.Time,
			UpdatedAt:   row.UpdatedAt.Time,
		}
		list = append(list, board)
	}
	return list, nil
}

// ListOwnedBoardAndUsers returns a list of boards that a user owns along with each board's associated members
// The SQL query uses a left join so it is possible that a board can have nullable board users.
func (r *repository) ListOwnedBoardAndUsers(ctx context.Context, userID uuid.UUID) ([]BoardAndUser, error) {
	rows, err := r.q.ListOwnedBoardAndUsers(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("repository: failed to list boards by user ID: %w", err)
	}

	// Convert database types into domain types
	list := []BoardAndUser{}
	for _, row := range rows {
		item := BoardAndUser{
			Board:           toBoard(row.Board),
			User:            toUser(row.User),
			BoardMembership: toBoardMembership(row.BoardMembership),
		}
		if err != nil {
			return nil, fmt.Errorf("repository: failed to transform db row to domain model: %w", err)
		}
		list = append(list, item)
	}
	return list, nil
}

// ListSharedBoardAndUsers returns a list of boards that a user belongs to along with a list of its associated members
// The SQL query uses a left join so it is possible that a board can have nullable board users.
func (r *repository) ListSharedBoardAndUsers(ctx context.Context, userID uuid.UUID) ([]BoardAndUser, error) {
	rows, err := r.q.ListSharedBoardAndUsers(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("repository: failed to list boards by user ID: %w", err)
	}

	// Convert database types into domain types
	list := []BoardAndUser{}
	for _, row := range rows {
		item := BoardAndUser{
			Board:           toBoard(row.Board),
			User:            toUser(row.User),
			BoardMembership: toBoardMembership(row.BoardMembership),
		}
		if err != nil {
			return nil, fmt.Errorf("repository: failed to transform db row to domain model: %w", err)
		}
		list = append(list, item)
	}
	return list, nil
}

// ListInvitesByBoard returns a list of board invites for a given board.
func (r *repository) ListInvitesByBoard(ctx context.Context, boardID uuid.UUID, status string) ([]models.Invite, error) {
	arg := db.ListInvitesByBoardParams{
		BoardID: pgtype.UUID{Bytes: boardID, Valid: true},
	}
	if status != "" {
		arg.Status = pgtype.Text{String: string(status), Valid: true}
	}
	rows, err := r.q.ListInvitesByBoard(ctx, arg)
	if err != nil {
		return []models.Invite{}, fmt.Errorf("repository: failed to list board invites: %w", err)
	}
	invites := []models.Invite{}
	for _, row := range rows {
		invite := toInvite(row)
		invites = append(invites, invite)
	}
	return invites, nil
}

// ListInvitesByReceiver returns a list of board invites along with board and sender details.
func (r *repository) ListInvitesByReceiver(ctx context.Context, receiverID uuid.UUID, status string) ([]InviteBoardSender, error) {
	arg := db.ListInvitesByReceiverParams{ReceiverID: pgtype.UUID{Bytes: receiverID, Valid: true}}
	if status != "" {
		arg.Status = pgtype.Text{String: string(status), Valid: true}
	}
	rows, err := r.q.ListInvitesByReceiver(ctx, arg)
	if err != nil {
		return []InviteBoardSender{}, fmt.Errorf("repository: failed to list board invites: %w", err)
	}
	inviteBoardSenders := []InviteBoardSender{}
	for _, row := range rows {
		invite := toInvite(row.BoardInvite)
		board := toBoard(row.Board)
		user := toUser(row.User)
		inviteBoardSenders = append(inviteBoardSenders, InviteBoardSender{invite, board, user})
	}
	return inviteBoardSenders, nil
}

// UpdateInvite updates an invite.
func (r *repository) UpdateInvite(ctx context.Context, invite models.Invite) error {
	arg := db.UpdateInviteParams(toInviteDB(invite))
	err := r.q.UpdateInvite(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInviteDoesNotExist
		}
		return fmt.Errorf("repository: failed to update invite: %w", err)
	}
	return nil
}

// DeleteBoard deletes a single board.
func (r *repository) DeleteBoard(ctx context.Context, boardID uuid.UUID) error {
	err := r.q.DeleteBoard(ctx, pgtype.UUID{Bytes: boardID, Valid: true})
	if err != nil {
		return fmt.Errorf("repository: failed to delete board: %w", err)
	}
	return nil
}

// BoardAndUser is a custom domain type to guard against nullable BoardMembership and
// User values in the case of left joins.
type BoardAndUser struct {
	Board           models.Board
	BoardMembership models.BoardMembership
	User            models.User
}

func toBoard(dbBoard db.Board) models.Board {
	return models.Board{
		ID:          dbBoard.ID.Bytes,
		Name:        &dbBoard.Name.String,
		Description: &dbBoard.Description.String,
		UserID:      dbBoard.UserID.Bytes,
		CreatedAt:   dbBoard.CreatedAt.Time,
		UpdatedAt:   dbBoard.UpdatedAt.Time,
	}
}

func toUser(dbUser db.User) models.User {
	return models.User{
		ID:        dbUser.ID.Bytes,
		Name:      dbUser.Name.String,
		Email:     &dbUser.Email.String,
		Password:  &dbUser.Password.String,
		IsGuest:   dbUser.IsGuest.Bool,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}
}

func toBoardMembership(dbBoardMembership db.BoardMembership) models.BoardMembership {
	return models.BoardMembership{
		ID:        dbBoardMembership.ID.Bytes,
		BoardID:   dbBoardMembership.UserID.Bytes,
		UserID:    dbBoardMembership.UserID.Bytes,
		Role:      models.BoardMembershipRole(dbBoardMembership.Role.String),
		CreatedAt: dbBoardMembership.CreatedAt.Time,
		UpdatedAt: dbBoardMembership.UpdatedAt.Time,
	}
}

type InviteBoardSender struct {
	Invite models.Invite
	Board  models.Board
	Sender models.User
}

func toInvite(row db.BoardInvite) models.Invite {
	return models.Invite{
		ID:         row.ID.Bytes,
		BoardID:    row.BoardID.Bytes,
		SenderID:   row.SenderID.Bytes,
		ReceiverID: row.ReceiverID.Bytes,
		Status:     models.InviteStatus(row.Status.String),
		CreatedAt:  row.CreatedAt.Time,
		UpdatedAt:  row.UpdatedAt.Time,
	}
}

func toInviteDB(invite models.Invite) db.BoardInvite {
	return db.BoardInvite{
		ID:         pgtype.UUID{Bytes: invite.ID, Valid: true},
		BoardID:    pgtype.UUID{Bytes: invite.BoardID, Valid: true},
		SenderID:   pgtype.UUID{Bytes: invite.SenderID, Valid: true},
		ReceiverID: pgtype.UUID{Bytes: invite.ReceiverID, Valid: true},
		Status:     pgtype.Text{String: string(invite.Status), Valid: true},
		CreatedAt:  pgtype.Timestamp{Time: invite.CreatedAt, Valid: true},
		UpdatedAt:  pgtype.Timestamp{Time: invite.UpdatedAt, Valid: true},
	}
}

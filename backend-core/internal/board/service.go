package board

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Wave-95/boards/backend-core/internal/amqp"
	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/pkg/logger"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/Wave-95/boards/shared/queues"
	"github.com/Wave-95/boards/shared/tasks"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	errInvalidID               = errors.New("ID not in UUID format")
	errBoardNotFound           = errors.New("Board not found")
	errInviteNotFound          = errors.New("Invite not found")
	errUnauthorized            = errors.New("User is not authorized")
	errUnsupportedInviteUpdate = errors.New("Invite update status is not supported")
	errInviteCancelled         = errors.New("Invite has been cancelled")
	errInvalidStatusFilter     = errors.New("Invalid status filter")
	defaultBoardDescription    = "My default board description"
)

// Service is an interface that represents all the capabilities of a board service.
type Service interface {
	CreateBoard(ctx context.Context, input CreateBoardInput) (models.Board, error)
	CreateInvites(ctx context.Context, input CreateInvitesInput) ([]models.Invite, error)

	GetBoard(ctx context.Context, boardID string) (models.Board, error)
	GetBoardWithMembers(ctx context.Context, boardID string) (BoardWithMembersDTO, error)
	GetInvite(ctx context.Context, inviteID string) (models.Invite, error)

	ListOwnedBoardsWithMembers(ctx context.Context, userID string) ([]BoardWithMembersDTO, error)
	ListSharedBoardsWithMembers(ctx context.Context, userID string) ([]BoardWithMembersDTO, error)
	ListInvitesByBoard(ctx context.Context, input ListInvitesByBoardInput) ([]InviteWithReceiverDTO, error)
	ListInvitesByReceiver(ctx context.Context, input ListInvitesByReceiverInput) ([]InviteWithBoardAndSenderDTO, error)

	UpdateInvite(ctx context.Context, input UpdateInviteInput) error
}

type service struct {
	repo      Repository
	amqp      amqp.Amqp
	validator validator.Validate
}

// NewService returns a service struct that implements the board service interface.
func NewService(repo Repository, amqp amqp.Amqp, validator validator.Validate) *service {
	return &service{
		repo:      repo,
		amqp:      amqp,
		validator: validator,
	}
}

// CreateBoard creates a new board and inserts the owner as the first member to that board. It will
// set the provided name and description or use defaults if none are provided.
func (s *service) CreateBoard(ctx context.Context, input CreateBoardInput) (models.Board, error) {
	// Validate input
	if err := s.validator.Struct(input); err != nil {
		return models.Board{}, err
	}
	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return models.Board{}, errInvalidID
	}

	// Create board name if none provided
	if input.Name == nil {
		boards, err := s.repo.ListOwnedBoards(ctx, userID)
		if err != nil {
			return models.Board{}, fmt.Errorf("service: failed to get existing boards when creating board: %w", err)
		}
		numBoards := len(boards)
		boardName := fmt.Sprintf("Board #%d", numBoards+1)
		input.Name = &boardName
	}

	// Use default board description if none provided
	if input.Description == nil {
		input.Description = &defaultBoardDescription
	}

	// Create new board
	boardID := uuid.New()
	now := time.Now()
	board := models.Board{
		ID:          boardID,
		Name:        input.Name,
		Description: input.Description,
		UserID:      userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.repo.CreateBoard(ctx, board); err != nil {
		return models.Board{}, fmt.Errorf("service: failed to create board: %w", err)
	}

	// Create membership for owner
	membershipID := uuid.New()
	membership := models.BoardMembership{
		ID:        membershipID,
		BoardID:   boardID,
		UserID:    userID,
		Role:      models.RoleAdmin,
		CreatedAt: now,
		UpdatedAt: now}
	if err := s.repo.CreateMembership(ctx, membership); err != nil {
		return models.Board{}, fmt.Errorf("service: failed to create board membership during create board: %w", err)
	}
	return board, nil
}

// CreateInvites creates board invites
func (s *service) CreateInvites(ctx context.Context, input CreateInvitesInput) ([]models.Invite, error) {
	boardID := input.BoardID
	senderID := input.SenderID
	invites := input.Invites
	// Parse IDs into UUIDs
	boardUUID, err := uuid.Parse(boardID)
	if err != nil {
		return nil, errInvalidID
	}
	senderUUID, err := uuid.Parse(senderID)
	if err != nil {
		return nil, errInvalidID
	}
	receiverUUIDs := []uuid.UUID{}
	for _, invite := range invites {
		receiverUUID, err := uuid.Parse(invite.ReceiverID)
		if err != nil {
			return nil, errInvalidID
		}
		receiverUUIDs = append(receiverUUIDs, receiverUUID)
	}

	// Check if user is authorized
	boardWithMembers, err := s.GetBoardWithMembers(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get board when creating board invites: %w", err)
	}
	if !UserIsAdmin(boardWithMembers, senderID) {
		return nil, errUnauthorized
	}

	// List existing pending invites
	listInvitesByBoardInput := ListInvitesByBoardInput{BoardID: boardID, UserID: senderID, Status: string(models.InviteStatusPending)}
	pendingInvites, err := s.ListInvitesByBoard(ctx, listInvitesByBoardInput)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get pending invites: %w", err)
	}

	invitesToInsert := []models.Invite{}
	now := time.Now()
	// Prepare invites to insert
	for _, receiverUUID := range receiverUUIDs {
		// If invite already exists, update the updated_at timestamp
		if existingInvite, ok := hasPendingInvite(receiverUUID, pendingInvites); ok {
			existingInvite.UpdatedAt = now
			invitesToInsert = append(invitesToInsert, existingInvite)
			continue
		}
		// Build new invite
		invite := models.Invite{
			ID:         uuid.New(),
			BoardID:    boardUUID,
			SenderID:   senderUUID,
			ReceiverID: receiverUUID,
			Status:     models.InviteStatusPending,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		invitesToInsert = append(invitesToInsert, invite)
	}

	err = s.repo.CreateInvites(ctx, invitesToInsert)
	if err != nil {
		return nil, fmt.Errorf("service: failed to create board invites: %w", err)
	}

	for _, invite := range invitesToInsert {
		s.amqp.Publish(queues.Notification, tasks.EmailInvite, invite)
	}

	return invitesToInsert, nil
}

// GetBoard returns a single board for a given board ID
func (s *service) GetBoard(ctx context.Context, boardID string) (models.Board, error) {
	boardUUID, err := uuid.Parse(boardID)
	if err != nil {
		return models.Board{}, fmt.Errorf("service: issue parsing boardID into UUID: %w", err)
	}
	board, err := s.repo.GetBoard(ctx, boardUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Board{}, errBoardNotFound
		}
		return models.Board{}, fmt.Errorf("service: failed to get board: %w", err)
	}
	return board, nil
}

// GetBoardWithMembers returns a single board with a list of associated members
func (s *service) GetBoardWithMembers(ctx context.Context, boardID string) (BoardWithMembersDTO, error) {
	logger := logger.FromContext(ctx)
	boardUUID, err := uuid.Parse(boardID)
	if err != nil {
		logger.Errorf("service: failed to parse boardID")
		return BoardWithMembersDTO{}, errInvalidID
	}
	rows, err := s.repo.GetBoardAndUsers(ctx, boardUUID)
	if err != nil {
		return BoardWithMembersDTO{}, fmt.Errorf("service: failed to get board with members: %w", err)
	}
	list := toBoardWithMembersDTO(rows)
	if len(list) == 0 {
		return BoardWithMembersDTO{}, errBoardNotFound
	}
	return list[0], nil
}

// GetInvite returns a single invite for a given invite ID
func (s *service) GetInvite(ctx context.Context, inviteID string) (models.Invite, error) {
	inviteUUID, err := uuid.Parse(inviteID)
	if err != nil {
		return models.Invite{}, fmt.Errorf("service: issue parsing inviteID into UUID: %w", err)
	}
	invite, err := s.repo.GetInvite(ctx, inviteUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Invite{}, errInviteNotFound
		}
		return models.Invite{}, fmt.Errorf("service: failed to get invite: %w", err)
	}
	return invite, nil
}

// ListOwnedBoardsWithMembers returns a list of boards that belong to a user
func (s *service) ListOwnedBoards(ctx context.Context, userID string) ([]models.Board, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("service: issue parsing userID into UUID: %w", err)
	}
	boards, err := s.repo.ListOwnedBoards(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get boards by user ID: %w", err)
	}
	return boards, nil
}

// ListOwnedBoardsWithMembers returns a list of boards that belong to a user along with a list of board members
func (s *service) ListOwnedBoardsWithMembers(ctx context.Context, userID string) ([]BoardWithMembersDTO, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("service: issue parsing userID into UUID: %w", err)
	}
	rows, err := s.repo.ListOwnedBoardAndUsers(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get boards by user ID: %w", err)
	}
	// transform rows into DTO
	list := toBoardWithMembersDTO(rows)
	return list, nil
}

// ListSharedBoardsWithMembers returns a list of boards that belong to a user along with a list of board members
func (s *service) ListSharedBoardsWithMembers(ctx context.Context, userID string) ([]BoardWithMembersDTO, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("service: issue parsing userID into UUID: %w", err)
	}
	rows, err := s.repo.ListSharedBoardAndUsers(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get boards by user ID: %w", err)
	}
	// transform rows into DTO
	list := toBoardWithMembersDTO(rows)
	return list, nil
}

// ListInvitesByBoard returns a list of invites belonging to a board. The user ID that is passed in should already
// be authenticated.
func (s *service) ListInvitesByBoard(ctx context.Context, input ListInvitesByBoardInput) ([]InviteWithReceiverDTO, error) {
	logger := logger.FromContext(ctx)
	boardID := input.BoardID
	status := input.Status
	userID := input.UserID
	boardUUID, err := uuid.Parse(boardID)
	if err != nil {
		return nil, errInvalidID
	}
	if status != "" {
		if !models.ValidInviteStatusFilter(status) {
			return []InviteWithReceiverDTO{}, errInvalidStatusFilter
		}
	}
	board, err := s.GetBoardWithMembers(ctx, boardID)
	if err != nil {
		return []InviteWithReceiverDTO{}, fmt.Errorf("service: failed to get board with members: %w", err)
	}
	if !UserHasAccess(board, userID) {
		return []InviteWithReceiverDTO{}, errUnauthorized
	}
	rows, err := s.repo.ListInvitesByBoard(ctx, boardUUID, status)
	if err != nil {
		logger.Errorf("service: failed to list invites by board: %w", err)
		return []InviteWithReceiverDTO{}, err
	}
	return toInviteWithReceiverDTO(rows), nil
}

// ListInvitesByReceiver returns a list of board invites for a given receiver. Each board invite element is augmented with
// sender and board details. The receiver ID should be the same as the authenticated user making the request.
func (s *service) ListInvitesByReceiver(ctx context.Context, input ListInvitesByReceiverInput) ([]InviteWithBoardAndSenderDTO, error) {
	receiverID := input.ReceiverID
	status := input.Status
	receiverUUID, err := uuid.Parse(receiverID)
	if err != nil {
		return nil, errInvalidID
	}
	if status != "" {
		if !models.ValidInviteStatusFilter(status) {
			return []InviteWithBoardAndSenderDTO{}, errInvalidStatusFilter
		}
	}
	inviteBoardSenders, err := s.repo.ListInvitesByReceiver(ctx, receiverUUID, status)
	if err != nil {
		return nil, fmt.Errorf("service: failed to list invites by receiver: %w", err)
	}
	dto := toInviteWithBoardAndSenderDTO(inviteBoardSenders)
	return dto, nil
}

// UpdateInvite updates a board invite. Only the sender of an invite can cancel the board invite, and only
// the receiver of an invite can accept or ignore the board invite.
func (s *service) UpdateInvite(ctx context.Context, input UpdateInviteInput) error {
	userUUID, err := uuid.Parse(input.UserID)
	if err != nil {
		return errInvalidID
	}
	invite, err := s.GetInvite(ctx, input.ID)
	if err != nil {
		return err
	}
	now := time.Now()
	switch input.Status {
	case string(models.InviteStatusAccepted):
		if invite.ReceiverID != userUUID {
			return errUnauthorized
		}
		if invite.Status == models.InviteStatusCancelled {
			return errInviteCancelled
		}
		// Add user to board
		membership := models.BoardMembership{
			ID:        uuid.New(),
			BoardID:   invite.BoardID,
			UserID:    invite.ReceiverID,
			Role:      models.RoleMember,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := s.repo.CreateMembership(ctx, membership); err != nil {
			return fmt.Errorf("service: failed to create membership when accepting invite: %w", err)
		}
	case string(models.InviteStatusIgnored):
		if invite.ReceiverID != userUUID {
			return errUnauthorized
		}
		if invite.Status == models.InviteStatusCancelled {
			return errInviteCancelled
		}
	case string(models.InviteStatusCancelled):
		if invite.SenderID != userUUID {
			return errUnauthorized
		}
	default:
		return errUnsupportedInviteUpdate
	}
	invite.UpdatedAt = now
	invite.Status = models.InviteStatus(input.Status)
	return s.repo.UpdateInvite(ctx, invite)
}

// toBoardWithMembersDTO transforms the BoardAndUser rows into a nested DTO struct
func toBoardWithMembersDTO(rows []BoardMembershipUser) []BoardWithMembersDTO {
	nestedList := []BoardWithMembersDTO{}
	boardIndex := make(map[uuid.UUID]int)
	for _, row := range rows {
		// If board does not exist, add it to the list
		if _, exists := boardIndex[row.Board.ID]; !exists {
			// Convert board domain model to DTO
			newItem := BoardWithMembersDTO{
				ID:          row.Board.ID,
				Name:        row.Board.Name,
				Description: row.Board.Description,
				UserID:      row.Board.UserID,
				Members:     []MemberDTO{},
				CreatedAt:   row.Board.CreatedAt,
				UpdatedAt:   row.Board.UpdatedAt,
			}
			boardIndex[row.Board.ID] = len(nestedList)
			nestedList = append(nestedList, newItem)
		}
		// Nest member/user details into board
		member := MemberDTO{
			ID:    row.User.ID,
			Name:  row.User.Name,
			Email: row.User.Email,
			Membership: MembershipDTO{
				Role:      string(row.Membership.Role),
				CreatedAt: row.Membership.CreatedAt,
				UpdatedAt: row.Membership.UpdatedAt,
			},
			CreatedAt: row.User.CreatedAt,
			UpdatedAt: row.User.UpdatedAt,
		}
		sliceIndex := boardIndex[row.Board.ID]
		board := nestedList[sliceIndex]
		board.Members = append(board.Members, member)
		nestedList[sliceIndex] = board
	}
	return nestedList
}

// hasPendingInvite checks to see if the receiver already has a pending invite, and returns that invite with a bool true.
func hasPendingInvite(receiverID uuid.UUID, pendingInvites []InviteWithReceiverDTO) (models.Invite, bool) {
	for _, pendingInvite := range pendingInvites {
		if pendingInvite.Receiver.ID == receiverID {
			invite := models.Invite{
				ID:         pendingInvite.ID,
				BoardID:    pendingInvite.BoardID,
				SenderID:   pendingInvite.SenderID,
				ReceiverID: pendingInvite.Receiver.ID,
				Status:     models.InviteStatusPending,
				CreatedAt:  pendingInvite.CreatedAt,
				UpdatedAt:  pendingInvite.UpdatedAt,
			}
			return invite, true
		}
	}
	return models.Invite{}, false
}

// toInviteWithBoardAndSenderDTO takes the flat structure from InviteBoardSender and maps it to the nested
// DTO structure.
func toInviteWithBoardAndSenderDTO(rows []InviteBoardSender) []InviteWithBoardAndSenderDTO {
	dto := []InviteWithBoardAndSenderDTO{}
	for _, row := range rows {
		row.Sender.Password = nil
		mappedRow := InviteWithBoardAndSenderDTO{
			ID:         row.Invite.ID,
			Board:      row.Board,
			Sender:     row.Sender,
			ReceiverID: row.Invite.ReceiverID,
			Status:     string(row.Invite.Status),
			CreatedAt:  row.Invite.CreatedAt,
			UpdatedAt:  row.Invite.UpdatedAt,
		}
		dto = append(dto, mappedRow)
	}
	return dto
}

// toInviteWithReceiverDTO takes the flat structure from InviteReceiver and maps it to the nested
// DTO structure.
func toInviteWithReceiverDTO(rows []InviteReceiver) []InviteWithReceiverDTO {
	dto := []InviteWithReceiverDTO{}
	for _, row := range rows {
		row.Receiver.Password = nil
		mappedRow := InviteWithReceiverDTO{
			ID:        row.Invite.ID,
			BoardID:   row.Invite.BoardID,
			SenderID:  row.Invite.SenderID,
			Receiver:  row.Receiver,
			Status:    string(row.Invite.Status),
			CreatedAt: row.Invite.CreatedAt,
			UpdatedAt: row.Invite.UpdatedAt,
		}
		dto = append(dto, mappedRow)
	}
	return dto
}

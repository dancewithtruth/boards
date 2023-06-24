// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: query.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createBoard = `-- name: CreateBoard :exec
INSERT INTO boards 
(id, name, description, user_id, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6)
`

type CreateBoardParams struct {
	ID          pgtype.UUID
	Name        pgtype.Text
	Description pgtype.Text
	UserID      pgtype.UUID
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
}

func (q *Queries) CreateBoard(ctx context.Context, arg CreateBoardParams) error {
	_, err := q.db.Exec(ctx, createBoard,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.UserID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const createInvite = `-- name: CreateInvite :exec
INSERT INTO board_invites
(id, board_id, sender_id, receiver_id, status, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
`

type CreateInviteParams struct {
	ID         pgtype.UUID
	BoardID    pgtype.UUID
	SenderID   pgtype.UUID
	ReceiverID pgtype.UUID
	Status     pgtype.Text
	CreatedAt  pgtype.Timestamp
	UpdatedAt  pgtype.Timestamp
}

func (q *Queries) CreateInvite(ctx context.Context, arg CreateInviteParams) error {
	_, err := q.db.Exec(ctx, createInvite,
		arg.ID,
		arg.BoardID,
		arg.SenderID,
		arg.ReceiverID,
		arg.Status,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const createMembership = `-- name: CreateMembership :exec
INSERT INTO board_memberships 
(id, user_id, board_id, role, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6)
`

type CreateMembershipParams struct {
	ID        pgtype.UUID
	UserID    pgtype.UUID
	BoardID   pgtype.UUID
	Role      pgtype.Text
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

func (q *Queries) CreateMembership(ctx context.Context, arg CreateMembershipParams) error {
	_, err := q.db.Exec(ctx, createMembership,
		arg.ID,
		arg.UserID,
		arg.BoardID,
		arg.Role,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const createPost = `-- name: CreatePost :exec
INSERT INTO posts
(id, board_id, user_id, content, pos_x, pos_y, color, height, z_index, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`

type CreatePostParams struct {
	ID        pgtype.UUID
	BoardID   pgtype.UUID
	UserID    pgtype.UUID
	Content   pgtype.Text
	PosX      pgtype.Int4
	PosY      pgtype.Int4
	Color     pgtype.Text
	Height    pgtype.Int4
	ZIndex    pgtype.Int4
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) error {
	_, err := q.db.Exec(ctx, createPost,
		arg.ID,
		arg.BoardID,
		arg.UserID,
		arg.Content,
		arg.PosX,
		arg.PosY,
		arg.Color,
		arg.Height,
		arg.ZIndex,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const deleteBoard = `-- name: DeleteBoard :exec
DELETE from boards WHERE id = $1
`

func (q *Queries) DeleteBoard(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteBoard, id)
	return err
}

const deletePost = `-- name: DeletePost :exec
DELETE from posts WHERE id = $1
`

func (q *Queries) DeletePost(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deletePost, id)
	return err
}

const getBoard = `-- name: GetBoard :one
SELECT id, name, description, user_id, created_at, updated_at FROM boards
WHERE boards.id = $1
`

func (q *Queries) GetBoard(ctx context.Context, id pgtype.UUID) (Board, error) {
	row := q.db.QueryRow(ctx, getBoard, id)
	var i Board
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getBoardAndUsers = `-- name: GetBoardAndUsers :many
SELECT boards.id, boards.name, boards.description, boards.user_id, boards.created_at, boards.updated_at, users.id, users.name, users.email, users.password, users.is_guest, users.created_at, users.updated_at, board_memberships.id, board_memberships.user_id, board_memberships.board_id, board_memberships.role, board_memberships.created_at, board_memberships.updated_at FROM boards
INNER JOIN board_memberships on board_memberships.board_id = boards.id
INNER JOIN users on board_memberships.user_id = users.id
WHERE boards.id = $1
ORDER BY boards.created_at DESC
`

type GetBoardAndUsersRow struct {
	Board           Board
	User            User
	BoardMembership BoardMembership
}

func (q *Queries) GetBoardAndUsers(ctx context.Context, id pgtype.UUID) ([]GetBoardAndUsersRow, error) {
	rows, err := q.db.Query(ctx, getBoardAndUsers, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBoardAndUsersRow
	for rows.Next() {
		var i GetBoardAndUsersRow
		if err := rows.Scan(
			&i.Board.ID,
			&i.Board.Name,
			&i.Board.Description,
			&i.Board.UserID,
			&i.Board.CreatedAt,
			&i.Board.UpdatedAt,
			&i.User.ID,
			&i.User.Name,
			&i.User.Email,
			&i.User.Password,
			&i.User.IsGuest,
			&i.User.CreatedAt,
			&i.User.UpdatedAt,
			&i.BoardMembership.ID,
			&i.BoardMembership.UserID,
			&i.BoardMembership.BoardID,
			&i.BoardMembership.Role,
			&i.BoardMembership.CreatedAt,
			&i.BoardMembership.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getInvite = `-- name: GetInvite :one
SELECT id, board_id, sender_id, receiver_id, status, created_at, updated_at FROM board_invites
WHERE board_invites.id = $1
`

func (q *Queries) GetInvite(ctx context.Context, id pgtype.UUID) (BoardInvite, error) {
	row := q.db.QueryRow(ctx, getInvite, id)
	var i BoardInvite
	err := row.Scan(
		&i.ID,
		&i.BoardID,
		&i.SenderID,
		&i.ReceiverID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPost = `-- name: GetPost :one
SELECT id, board_id, user_id, content, pos_x, pos_y, color, height, z_index, created_at, updated_at FROM posts
WHERE posts.id = $1
`

func (q *Queries) GetPost(ctx context.Context, id pgtype.UUID) (Post, error) {
	row := q.db.QueryRow(ctx, getPost, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.BoardID,
		&i.UserID,
		&i.Content,
		&i.PosX,
		&i.PosY,
		&i.Color,
		&i.Height,
		&i.ZIndex,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, email, password, is_guest, created_at, updated_at FROM users
WHERE users.email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email pgtype.Text) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.IsGuest,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listInvitesByBoard = `-- name: ListInvitesByBoard :many
SELECT board_invites.id, board_invites.board_id, board_invites.sender_id, board_invites.receiver_id, board_invites.status, board_invites.created_at, board_invites.updated_at, users.id, users.name, users.email, users.password, users.is_guest, users.created_at, users.updated_at FROM board_invites
INNER JOIN users on users.id = board_invites.receiver_id
WHERE board_invites.board_id = $1 AND
(status = $2 OR $2 IS NULL)
ORDER BY board_invites.updated_at DESC
`

type ListInvitesByBoardParams struct {
	BoardID pgtype.UUID
	Status  pgtype.Text
}

type ListInvitesByBoardRow struct {
	BoardInvite BoardInvite
	User        User
}

func (q *Queries) ListInvitesByBoard(ctx context.Context, arg ListInvitesByBoardParams) ([]ListInvitesByBoardRow, error) {
	rows, err := q.db.Query(ctx, listInvitesByBoard, arg.BoardID, arg.Status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListInvitesByBoardRow
	for rows.Next() {
		var i ListInvitesByBoardRow
		if err := rows.Scan(
			&i.BoardInvite.ID,
			&i.BoardInvite.BoardID,
			&i.BoardInvite.SenderID,
			&i.BoardInvite.ReceiverID,
			&i.BoardInvite.Status,
			&i.BoardInvite.CreatedAt,
			&i.BoardInvite.UpdatedAt,
			&i.User.ID,
			&i.User.Name,
			&i.User.Email,
			&i.User.Password,
			&i.User.IsGuest,
			&i.User.CreatedAt,
			&i.User.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listInvitesByReceiver = `-- name: ListInvitesByReceiver :many
SELECT board_invites.id, board_invites.board_id, board_invites.sender_id, board_invites.receiver_id, board_invites.status, board_invites.created_at, board_invites.updated_at, users.id, users.name, users.email, users.password, users.is_guest, users.created_at, users.updated_at, boards.id, boards.name, boards.description, boards.user_id, boards.created_at, boards.updated_at FROM board_invites
INNER JOIN boards on boards.id = board_invites.board_id
INNER JOIN users on users.id = board_invites.sender_id 
WHERE board_invites.receiver_id = $1 AND
(status = $2 OR $2 IS NULL)
ORDER BY board_invites.updated_at DESC
`

type ListInvitesByReceiverParams struct {
	ReceiverID pgtype.UUID
	Status     pgtype.Text
}

type ListInvitesByReceiverRow struct {
	BoardInvite BoardInvite
	User        User
	Board       Board
}

func (q *Queries) ListInvitesByReceiver(ctx context.Context, arg ListInvitesByReceiverParams) ([]ListInvitesByReceiverRow, error) {
	rows, err := q.db.Query(ctx, listInvitesByReceiver, arg.ReceiverID, arg.Status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListInvitesByReceiverRow
	for rows.Next() {
		var i ListInvitesByReceiverRow
		if err := rows.Scan(
			&i.BoardInvite.ID,
			&i.BoardInvite.BoardID,
			&i.BoardInvite.SenderID,
			&i.BoardInvite.ReceiverID,
			&i.BoardInvite.Status,
			&i.BoardInvite.CreatedAt,
			&i.BoardInvite.UpdatedAt,
			&i.User.ID,
			&i.User.Name,
			&i.User.Email,
			&i.User.Password,
			&i.User.IsGuest,
			&i.User.CreatedAt,
			&i.User.UpdatedAt,
			&i.Board.ID,
			&i.Board.Name,
			&i.Board.Description,
			&i.Board.UserID,
			&i.Board.CreatedAt,
			&i.Board.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listOwnedBoardAndUsers = `-- name: ListOwnedBoardAndUsers :many
SELECT boards.id, boards.name, boards.description, boards.user_id, boards.created_at, boards.updated_at, users.id, users.name, users.email, users.password, users.is_guest, users.created_at, users.updated_at, board_memberships.id, board_memberships.user_id, board_memberships.board_id, board_memberships.role, board_memberships.created_at, board_memberships.updated_at FROM boards
INNER JOIN board_memberships on board_memberships.board_id = boards.id
INNER JOIN users on board_memberships.user_id = users.id
WHERE boards.user_id = $1
ORDER BY boards.created_at DESC
`

type ListOwnedBoardAndUsersRow struct {
	Board           Board
	User            User
	BoardMembership BoardMembership
}

func (q *Queries) ListOwnedBoardAndUsers(ctx context.Context, userID pgtype.UUID) ([]ListOwnedBoardAndUsersRow, error) {
	rows, err := q.db.Query(ctx, listOwnedBoardAndUsers, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListOwnedBoardAndUsersRow
	for rows.Next() {
		var i ListOwnedBoardAndUsersRow
		if err := rows.Scan(
			&i.Board.ID,
			&i.Board.Name,
			&i.Board.Description,
			&i.Board.UserID,
			&i.Board.CreatedAt,
			&i.Board.UpdatedAt,
			&i.User.ID,
			&i.User.Name,
			&i.User.Email,
			&i.User.Password,
			&i.User.IsGuest,
			&i.User.CreatedAt,
			&i.User.UpdatedAt,
			&i.BoardMembership.ID,
			&i.BoardMembership.UserID,
			&i.BoardMembership.BoardID,
			&i.BoardMembership.Role,
			&i.BoardMembership.CreatedAt,
			&i.BoardMembership.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listOwnedBoards = `-- name: ListOwnedBoards :many
SELECT id, name, description, user_id, created_at, updated_at FROM boards
WHERE boards.user_id = $1
ORDER BY boards.created_at DESC
`

func (q *Queries) ListOwnedBoards(ctx context.Context, userID pgtype.UUID) ([]Board, error) {
	rows, err := q.db.Query(ctx, listOwnedBoards, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Board
	for rows.Next() {
		var i Board
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPosts = `-- name: ListPosts :many
SELECT id, board_id, user_id, content, pos_x, pos_y, color, height, z_index, created_at, updated_at FROM posts
WHERE posts.board_id = $1
`

func (q *Queries) ListPosts(ctx context.Context, boardID pgtype.UUID) ([]Post, error) {
	rows, err := q.db.Query(ctx, listPosts, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.BoardID,
			&i.UserID,
			&i.Content,
			&i.PosX,
			&i.PosY,
			&i.Color,
			&i.Height,
			&i.ZIndex,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSharedBoardAndUsers = `-- name: ListSharedBoardAndUsers :many
SELECT boards.id, boards.name, boards.description, boards.user_id, boards.created_at, boards.updated_at, users.id, users.name, users.email, users.password, users.is_guest, users.created_at, users.updated_at, board_memberships.id, board_memberships.user_id, board_memberships.board_id, board_memberships.role, board_memberships.created_at, board_memberships.updated_at FROM boards
INNER JOIN board_memberships on board_memberships.board_id = boards.id
INNER JOIN users on board_memberships.user_id = users.id
WHERE board_memberships.user_id = $1
AND board_memberships.role = 'MEMBER'
ORDER BY board_memberships.created_at DESC
`

type ListSharedBoardAndUsersRow struct {
	Board           Board
	User            User
	BoardMembership BoardMembership
}

func (q *Queries) ListSharedBoardAndUsers(ctx context.Context, userID pgtype.UUID) ([]ListSharedBoardAndUsersRow, error) {
	rows, err := q.db.Query(ctx, listSharedBoardAndUsers, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListSharedBoardAndUsersRow
	for rows.Next() {
		var i ListSharedBoardAndUsersRow
		if err := rows.Scan(
			&i.Board.ID,
			&i.Board.Name,
			&i.Board.Description,
			&i.Board.UserID,
			&i.Board.CreatedAt,
			&i.Board.UpdatedAt,
			&i.User.ID,
			&i.User.Name,
			&i.User.Email,
			&i.User.Password,
			&i.User.IsGuest,
			&i.User.CreatedAt,
			&i.User.UpdatedAt,
			&i.BoardMembership.ID,
			&i.BoardMembership.UserID,
			&i.BoardMembership.BoardID,
			&i.BoardMembership.Role,
			&i.BoardMembership.CreatedAt,
			&i.BoardMembership.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUsersByEmail = `-- name: ListUsersByEmail :many
SELECT id, name, email, password, is_guest, created_at, updated_at FROM users
ORDER BY levenshtein(users.email, $1) LIMIT 10
`

func (q *Queries) ListUsersByEmail(ctx context.Context, levenshtein interface{}) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsersByEmail, levenshtein)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Password,
			&i.IsGuest,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateInvite = `-- name: UpdateInvite :exec
UPDATE board_invites SET
(id, board_id, sender_id, receiver_id, status, created_at, updated_at) =
($1, $2, $3, $4, $5, $6, $7) WHERE id = $1
`

type UpdateInviteParams struct {
	ID         pgtype.UUID
	BoardID    pgtype.UUID
	SenderID   pgtype.UUID
	ReceiverID pgtype.UUID
	Status     pgtype.Text
	CreatedAt  pgtype.Timestamp
	UpdatedAt  pgtype.Timestamp
}

func (q *Queries) UpdateInvite(ctx context.Context, arg UpdateInviteParams) error {
	_, err := q.db.Exec(ctx, updateInvite,
		arg.ID,
		arg.BoardID,
		arg.SenderID,
		arg.ReceiverID,
		arg.Status,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const updatePost = `-- name: UpdatePost :exec
UPDATE posts SET
(id, board_id, user_id, content, pos_x, pos_y, color, height, z_index, created_at, updated_at) =
($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) WHERE id = $1
`

type UpdatePostParams struct {
	ID        pgtype.UUID
	BoardID   pgtype.UUID
	UserID    pgtype.UUID
	Content   pgtype.Text
	PosX      pgtype.Int4
	PosY      pgtype.Int4
	Color     pgtype.Text
	Height    pgtype.Int4
	ZIndex    pgtype.Int4
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

func (q *Queries) UpdatePost(ctx context.Context, arg UpdatePostParams) error {
	_, err := q.db.Exec(ctx, updatePost,
		arg.ID,
		arg.BoardID,
		arg.UserID,
		arg.Content,
		arg.PosX,
		arg.PosY,
		arg.Color,
		arg.Height,
		arg.ZIndex,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

-- name: CreateBoard :exec
INSERT INTO boards 
(id, name, description, user_id, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetBoard :one
SELECT * FROM boards
WHERE boards.id = $1;

-- name: GetBoardAndUsers :many
SELECT sqlc.embed(boards), sqlc.embed(users), sqlc.embed(board_memberships) FROM boards
LEFT JOIN board_memberships on board_memberships.board_id = boards.id
LEFT JOIN users on board_memberships.user_id = users.id
WHERE boards.id = $1;

-- name: ListOwnedBoards :many
SELECT * FROM boards
WHERE boards.user_id = $1;

-- name: ListOwnedBoardAndUsers :many
SELECT sqlc.embed(boards), sqlc.embed(users), sqlc.embed(board_memberships) FROM boards
LEFT JOIN board_memberships on board_memberships.board_id = boards.id
LEFT JOIN users on board_memberships.user_id = users.id
WHERE boards.user_id = $1;

-- name: CreateBoardInvite :exec
INSERT INTO board_invites
(id, user_id, board_id, status, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6);

-- name: CreateMembership :exec
INSERT INTO board_memberships 
(id, user_id, board_id, role, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6);

-- name: DeleteBoard :exec
DELETE from boards WHERE id = $1;


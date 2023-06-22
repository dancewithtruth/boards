-- name: CreateBoard :exec
INSERT INTO boards 
(id, name, description, user_id, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetBoard :one
SELECT * FROM boards
WHERE boards.id = $1;

-- name: GetBoardAndUsers :many
SELECT sqlc.embed(boards), sqlc.embed(users), sqlc.embed(board_memberships) FROM boards
INNER JOIN board_memberships on board_memberships.board_id = boards.id
INNER JOIN users on board_memberships.user_id = users.id
WHERE boards.id = $1
ORDER BY boards.created_at DESC;

-- name: ListOwnedBoards :many
SELECT * FROM boards
WHERE boards.user_id = $1
ORDER BY boards.created_at DESC;

-- name: ListOwnedBoardAndUsers :many
SELECT sqlc.embed(boards), sqlc.embed(users), sqlc.embed(board_memberships) FROM boards
INNER JOIN board_memberships on board_memberships.board_id = boards.id
INNER JOIN users on board_memberships.user_id = users.id
WHERE boards.user_id = $1
ORDER BY boards.created_at DESC;

-- name: ListSharedBoardAndUsers :many
SELECT sqlc.embed(boards), sqlc.embed(users), sqlc.embed(board_memberships) FROM boards
INNER JOIN board_memberships on board_memberships.board_id = boards.id
INNER JOIN users on board_memberships.user_id = users.id
WHERE board_memberships.user_id = $1
AND board_memberships.role = 'MEMBER'
ORDER BY board_memberships.created_at DESC;

-- name: CreateMembership :exec
INSERT INTO board_memberships 
(id, user_id, board_id, role, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6);

-- name: DeleteBoard :exec
DELETE from boards WHERE id = $1;

-- name: CreatePost :exec
INSERT INTO posts
(id, board_id, user_id, content, pos_x, pos_y, color, height, z_index, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: GetPost :one
SELECT * FROM posts
WHERE posts.id = $1;

-- name: ListPosts :many
SELECT * FROM posts
WHERE posts.board_id = $1;

-- name: UpdatePost :exec
UPDATE posts SET
(id, board_id, user_id, content, pos_x, pos_y, color, height, z_index, created_at, updated_at) =
($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) WHERE id = $1;

-- name: DeletePost :exec
DELETE from posts WHERE id = $1;

-- name: ListUsersByEmail :many
SELECT * FROM users
ORDER BY levenshtein(users.email, $1) LIMIT 10;

-- name: CreateInvite :exec
INSERT INTO board_invites
(id, board_id, sender_id, receiver_id, status, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: UpdateInvite :exec
UPDATE board_invites SET
(id, board_id, sender_id, receiver_id, status, created_at, updated_at) =
($1, $2, $3, $4, $5, $6, $7) WHERE id = $1;

-- name: ListInvitesByBoard :many
SELECT * FROM board_invites
WHERE board_invites.board_id = $1
ORDER BY board_invites.updated_at DESC;

-- name: ListInvitesByBoardFilterStatus :many
SELECT * FROM board_invites
WHERE board_invites.board_id = $1
AND board_invites.status = $2
ORDER BY board_invites.updated_at DESC;

-- name: ListInvitesByReceiver :many
SELECT * FROM board_invites
WHERE board_invites.receiver_id = $1
ORDER BY board_invites.updated_at DESC;

-- name: ListInvitesByReceiverFilterStatus :many
SELECT * FROM board_invites
WHERE board_invites.receiver_id = $1
AND board_invites.status = $2
ORDER BY board_invites.updated_at DESC;

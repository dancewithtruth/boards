BEGIN;

-- Drop the index on the "user_id" column
DROP INDEX IF EXISTS idx_boards_user_id;

-- Drop the "boards" table
DROP TABLE IF EXISTS boards;

COMMIT;
BEGIN;

-- Create the "boards" table
CREATE TABLE IF NOT EXISTS boards(
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    user_id UUID REFERENCES users(id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create an index on the "user_id" column
CREATE INDEX idx_boards_user_id ON boards (user_id);

COMMIT;
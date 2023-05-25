CREATE TABLE IF NOT EXISTS users_boards (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  board_id UUID REFERENCES boards(id) ON DELETE CASCADE,
  CONSTRAINT unique_user_board UNIQUE (user_id, board_id)
);

CREATE INDEX idx_user_id ON users_boards (user_id);
CREATE INDEX idx_board_id ON users_boards (board_id);

CREATE TABLE IF NOT EXISTS board_members (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  board_id UUID REFERENCES boards(id) ON DELETE CASCADE,
  role VARCHAR(20) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  CONSTRAINT unique_board_members UNIQUE (user_id, board_id)
);

CREATE INDEX idx_board_members_user_id ON board_members (user_id);
CREATE INDEX idx_board_members_board_id ON board_members (board_id);

CREATE TABLE IF NOT EXISTS board_memberships (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  board_id UUID REFERENCES boards(id) ON DELETE CASCADE,
  role VARCHAR(20) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  CONSTRAINT unique_board_memberships UNIQUE (user_id, board_id)
);

CREATE INDEX idx_board_memberships_user_id ON board_memberships (user_id);
CREATE INDEX idx_board_memberships_board_id ON board_memberships (board_id);

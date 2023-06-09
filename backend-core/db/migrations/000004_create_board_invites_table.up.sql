CREATE TABLE IF NOT EXISTS board_invites (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  board_id UUID REFERENCES boards(id) ON DELETE CASCADE,
  status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  CONSTRAINT unique_board_invites UNIQUE (user_id, board_id)
);

CREATE INDEX idx_board_invites_user_id ON board_invites (user_id);
CREATE INDEX idx_board_invites_board_id ON board_invites (board_id);

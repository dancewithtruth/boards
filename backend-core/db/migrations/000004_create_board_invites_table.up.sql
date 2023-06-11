CREATE TABLE IF NOT EXISTS board_invites (
  id UUID PRIMARY KEY,
  board_id UUID REFERENCES boards(id) ON DELETE CASCADE,
  sender_id UUID REFERENCES users(id) ON DELETE CASCADE,
  receiver_id UUID REFERENCES users(id) ON DELETE CASCADE,
  status VARCHAR(20) DEFAULT 'PENDING',
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_board_invites_user_id ON board_invites (receiver_id);
CREATE INDEX idx_board_invites_board_id ON board_invites (board_id);
CREATE EXTENSION IF NOT EXISTS pg_trgm;
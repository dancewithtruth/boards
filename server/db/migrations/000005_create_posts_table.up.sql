CREATE TABLE IF NOT EXISTS posts (
  id UUID PRIMARY KEY,
  board_id UUID REFERENCES boards(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  content TEXT,
  pos_x INTEGER,
  pos_y INTEGER,
  color VARCHAR(7),
  height DECIMAL,
  z_index INTEGER,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

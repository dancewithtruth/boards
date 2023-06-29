CREATE TABLE IF NOT EXISTS post_groups (
  id UUID PRIMARY KEY,
  title VARCHAR(50),
  pos_x INTEGER,
  pos_y INTEGER,
  z_index INTEGER,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

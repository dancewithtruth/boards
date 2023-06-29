-- Drop foreign key constraint and delete trigger
ALTER TABLE posts
DROP CONSTRAINT IF EXISTS fk_post_group;

ALTER TABLE posts
DROP COLUMN post_group_id,
DROP COLUMN post_order,
ADD COLUMN pos_x INTEGER,
ADD COLUMN pos_y INTEGER,
ADD COLUMN z_index INTEGER;

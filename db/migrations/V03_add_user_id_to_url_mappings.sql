-- +migrate Up
ALTER TABLE url_mappings
ADD COLUMN user_id TEXT REFERENCES users(id);

-- +migrate Down
ALTER TABLE url_mappings
DROP COLUMN user_id;
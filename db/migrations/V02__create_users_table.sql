-- +migrate Up
CREATE TABLE users (
  id TEXT PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  name TEXT,
  picture_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE users;
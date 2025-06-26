-- +migrate Up
CREATE TABLE url_mappings (
  key TEXT PRIMARY KEY,
  original_url TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_original_url ON url_mappings (original_url);

-- +migrate Down
DROP TABLE url_mappings;
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS posts (
  id          uuid PRIMARY KEY,
  title       varchar(255)       NOT NULL,
  body        text        NOT NULL,
  author_id   uuid        NOT NULL,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT posts_author_fk FOREIGN KEY (author_id)
    REFERENCES users(id) ON DELETE RESTRICT
);
CREATE INDEX IF NOT EXISTS posts_author_id_idx ON posts (author_id);
CREATE INDEX IF NOT EXISTS posts_created_at_idx ON posts (created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS posts_created_at_idx;
DROP INDEX IF EXISTS posts_author_id_idx;
DROP TABLE IF EXISTS posts;
-- +goose StatementEnd

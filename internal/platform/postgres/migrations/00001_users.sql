-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id            uuid PRIMARY KEY NOT NULL,
  email         varchar(255)     NOT NULL,
  password_hash varchar(255)     NOT NULL,
  created_at    timestamptz      NOT NULL DEFAULT now(),
  updated_at    timestamptz      NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX IF NOT EXISTS users_email_lower_unique ON users (lower(email));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS users_email_lower_unique;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE IF NOT EXISTS "users"
(
    id         UUID                 DEFAULT uuid_generate_v4() NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    email      TEXT        NOT NULL UNIQUE,
    password   TEXT        NOT NULL,
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "users";
-- +goose StatementEnd

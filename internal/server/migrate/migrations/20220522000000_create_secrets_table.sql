-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE IF NOT EXISTS "secrets"
(
    id         UUID                  DEFAULT uuid_generate_v4() NOT NULL UNIQUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    user_id    UUID         NOT NULL,
    type       VARCHAR(255) NOT NULL,
    name       VARCHAR(255) NOT NULL,
    content    BYTEA        NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
            REFERENCES users (id)
);
CREATE UNIQUE INDEX IF NOT EXISTS secrets_unique_user_id_name
    ON secrets (user_id, name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "users";
-- +goose StatementEnd

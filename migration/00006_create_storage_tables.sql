-- The client stores account/char specific data on the server.

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS account_storage (
    account_id      integer NOT NULL REFERENCES accounts ON DELETE CASCADE,
    updated_at      timestamp NOT NULL DEFAULT now(),
    type            integer NOT NULL,
    data            bytea, -- compressed binary blob
    PRIMARY KEY (account_id, type)
);

CREATE TABLE IF NOT EXISTS character_storage (
    character_id    integer NOT NULL REFERENCES characters ON DELETE CASCADE,
    updated_at      timestamp NOT NULL DEFAULT now(),
    type            integer NOT NULL,
    data            bytea, -- compressed binary blob
    PRIMARY KEY (character_id, type)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS account_storage;
DROP TABLE IF EXISTS character_storage;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS account_storage (
    account_id      integer NOT NULL REFERENCES accounts ON DELETE CASCADE,
    updated_at      timestamp NOT NULL DEFAULT now(),
    type            integer NOT NULL CHECK (type >= 0 AND type <= 7),
    data            bytea, -- compressed binary blob
    PRIMARY KEY (account_id, type)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS account_storage;
-- +goose StatementEnd

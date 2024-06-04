-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS realms (
    id              serial PRIMARY KEY,
    created_at      timestamp NOT NULL DEFAULT now(),
    name            varchar(100) NOT NULL,
    type            integer NOT NULL,
    host            varchar(100) NOT NULL,
    region          integer NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS realms_host_unique_idx ON realms (lower(host));
CREATE UNIQUE INDEX IF NOT EXISTS realms_name_unique_idx ON realms (lower(name));

CREATE TABLE IF NOT EXISTS accounts (
    id              serial PRIMARY KEY,
    created_at      timestamp NOT NULL DEFAULT now(),
    last_login      timestamp,
    username        varchar(16) NOT NULL,
    srp_verifier    varchar(64) NOT NULL, -- 32 byte hex string
    srp_salt        varchar(64) NOT NULL, -- 32 byte hex string
    email           varchar(100) NOT NULL,
    realm_id        integer NOT NULL REFERENCES realms ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS accounts_email_unique_idx ON accounts (lower(email));
CREATE UNIQUE INDEX IF NOT EXISTS accounts_username_unique_idx ON accounts (upper(username));

CREATE TABLE IF NOT EXISTS sessions (
    id              serial PRIMARY KEY,
    account         integer UNIQUE NOT NULL REFERENCES accounts ON DELETE CASCADE,
    session_key     varchar(80) NOT NULL, -- 40 byte hex string
    connected       integer NOT NULL,
    connected_at    timestamp,
    disconnected_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS realms;
-- +goose StatementEnd

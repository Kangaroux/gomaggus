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
CREATE TABLE IF NOT EXISTS accounts (
    id              serial PRIMARY KEY,
    created_at      timestamp NOT NULL DEFAULT now(),
    last_login      timestamp,
    username        varchar(16) NOT NULL,
    srp_verifier    varchar(64) NOT NULL, -- 32 byte hex string
    srp_salt        varchar(64) NOT NULL, -- 32 byte hex string
    email           varchar(100) NOT NULL,
    realm_id        integer NOT NULL REFERENCES realms(id)
);
CREATE TABLE IF NOT EXISTS sessions (
    id              serial PRIMARY KEY,
    connected       integer NOT NULL DEFAULT 0,
    session_key     varchar(80) NOT NULL, -- 40 byte hex string
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

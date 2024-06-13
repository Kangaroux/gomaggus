-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS characters (
    id            serial PRIMARY KEY,
    created_at    timestamp NOT NULL DEFAULT now(),
    last_login    timestamp,
    name          varchar(12) NOT NULL,
    account_id    integer NOT NULL REFERENCES accounts (id),
    realm_id      integer NOT NULL REFERENCES realms (id),
    race          smallint NOT NULL,
    class         smallint NOT NULL,
    gender        smallint NOT NULL,
    skin_color    smallint NOT NULL,
    face          smallint NOT NULL,
    hair_style    smallint NOT NULL,
    hair_color    smallint NOT NULL,
    facial_hair   smallint NOT NULL,
    outfit_id     smallint NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS characters_name_realm_unique_idx ON characters (lower(name), realm_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS characters;
-- +goose StatementEnd

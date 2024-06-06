-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS characters (
    id            serial PRIMARY KEY,
    created_at    timestamp NOT NULL DEFAULT now(),
    last_login    timestamp,
    name          varchar(12) NOT NULL,
    account_id    integer NOT NULL REFERENCES accounts (id),
    realm_id      integer NOT NULL REFERENCES realms (id),
    race          char NOT NULL,
    class         char NOT NULL,
    gender        char NOT NULL,
    skin_color    char NOT NULL,
    face          char NOT NULL,
    hair_style    char NOT NULL,
    hair_color    char NOT NULL,
    facial_hair   char NOT NULL,
    outfit_id     char NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS characters_name_realm_unique_idx ON characters (lower(name), realm_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS characters;
-- +goose StatementEnd

-- Removes `sessions.id` and makes the account_id the primary key.

-- +goose Up
-- +goose StatementBegin
TRUNCATE TABLE sessions;

-- Remove `id` column
ALTER TABLE sessions DROP CONSTRAINT sessions_pkey;
ALTER TABLE sessions DROP id;

-- Make account_id the new PK
ALTER TABLE sessions ADD PRIMARY KEY (account_id);

-- Remove redundant `UNIQUE` constraint on account_id (since it's a PK now)
ALTER TABLE sessions DROP CONSTRAINT sessions_account_id_key;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE sessions;

-- Change account_id to no longer be the PK
ALTER TABLE sessions DROP CONSTRAINT sessions_pkey;

-- Add `id` as the new PK
ALTER TABLE sessions ADD id serial PRIMARY KEY;
ALTER TABLE sessions ADD UNIQUE (account_id);
-- +goose StatementEnd

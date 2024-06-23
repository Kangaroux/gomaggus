-- The client wants to know what the size is but the server doesn't want to decompress to check

-- +goose Up
-- +goose StatementBegin

--
TRUNCATE TABLE account_storage;
ALTER TABLE account_storage ADD uncompressed_size INTEGER NOT NULL;
TRUNCATE TABLE character_storage;
ALTER TABLE character_storage ADD uncompressed_size INTEGER NOT NULL;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE account_storage DROP uncompressed_size;
ALTER TABLE character_storage DROP uncompressed_size;
-- +goose StatementEnd

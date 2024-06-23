-- The extra cosmetic can be more than just facial hair; just clarifying naming here.

-- +goose Up
-- +goose StatementBegin
ALTER TABLE characters RENAME facial_hair TO extra_cosmetic;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE characters RENAME extra_cosmetic TO facial_hair;
-- +goose StatementEnd

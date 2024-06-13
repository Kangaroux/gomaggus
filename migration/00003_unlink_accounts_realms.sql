-- +goose Up
-- +goose StatementBegin
ALTER TABLE accounts DROP realm_id;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- This column should be NOT NULL, but it needs to be populated first. Start by adding it as nullable
ALTER TABLE accounts ADD realm_id integer REFERENCES realms ON DELETE CASCADE;

-- Populate the column with any available realm ID
UPDATE accounts SET realm_id = (SELECT id FROM realms LIMIT 1);

-- Make it not nullable
ALTER TABLE accounts ALTER realm_id SET NOT NULL;

-- +goose StatementEnd

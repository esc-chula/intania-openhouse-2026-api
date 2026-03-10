-- +goose Up
-- +goose StatementBegin
CREATE TYPE booth_category AS ENUM ('Department', 'Club', 'Exhibition');
ALTER TABLE booths ADD COLUMN IF NOT EXISTS category booth_category NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE booths DROP COLUMN IF EXISTS category;
-- +goose StatementEnd

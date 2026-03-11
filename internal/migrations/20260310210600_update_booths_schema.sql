-- +goose Up
-- +goose StatementBegin
CREATE TYPE booth_category AS ENUM ('department', 'club', 'exhibition');
ALTER TABLE booths ADD COLUMN IF NOT EXISTS category booth_category NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE booths DROP COLUMN IF EXISTS category;
DROP TYPE IF EXISTS booth_category;
-- +goose StatementEnd

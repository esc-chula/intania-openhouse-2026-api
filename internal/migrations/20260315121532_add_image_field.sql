-- +goose Up
-- +goose StatementBegin
ALTER TABLE workshops ADD COLUMN IF NOT EXISTS image TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workshops DROP COLUMN IF EXISTS image;
-- +goose StatementEnd

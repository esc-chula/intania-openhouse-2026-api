-- +goose Up
-- +goose StatementBegin
ALTER TABLE activities ADD COLUMN event_date DATE;
UPDATE activities SET event_date = start_time::date WHERE event_date IS NULL;
ALTER TABLE activities ALTER COLUMN event_date SET NOT NULL;
ALTER TABLE activities ADD COLUMN link TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE activities DROP COLUMN IF EXISTS event_date;
ALTER TABLE activities DROP COLUMN IF EXISTS link;
-- +goose StatementEnd

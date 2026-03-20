-- +goose Up
-- +goose StatementBegin
ALTER TABLE activities
  ALTER COLUMN start_time TYPE TIME USING (start_time::time),
  ALTER COLUMN end_time TYPE TIME USING (end_time::time);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE activities
  ALTER COLUMN start_time TYPE TIMESTAMP WITH TIME ZONE USING ((event_date::timestamp + start_time::time)::timestamptz),
  ALTER COLUMN end_time TYPE TIMESTAMP WITH TIME ZONE USING ((event_date::timestamp + end_time::time)::timestamptz);
-- +goose StatementEnd

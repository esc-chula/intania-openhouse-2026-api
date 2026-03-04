-- +goose Up
-- +goose StatementBegin
ALTER TYPE booking_status ADD VALUE IF NOT EXISTS 'Attended';
ALTER TYPE booking_status ADD VALUE IF NOT EXISTS 'Absent';
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS checked_in_at TIMESTAMP WITH TIME ZONE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE bookings DROP COLUMN IF EXISTS checked_in_at;
-- +goose StatementEnd

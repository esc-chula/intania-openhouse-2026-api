-- +goose Up
-- +goose StatementBegin

-- Drop the old index that only checked for 'Confirmed'
DROP INDEX IF EXISTS uniq_user_workshop_confirmed;

-- Create the new partial unique index for active/terminal states
CREATE UNIQUE INDEX uniq_user_workshop_active_booking 
ON bookings (user_id, workshop_id) 
WHERE status IN ('Confirmed'::booking_status, 'Attended'::booking_status, 'Absent'::booking_status);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS uniq_user_workshop_active_booking;

CREATE UNIQUE INDEX uniq_user_workshop_confirmed 
ON bookings (user_id, workshop_id) 
WHERE status = 'Confirmed'::booking_status;
-- +goose StatementEnd

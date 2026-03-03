-- +goose Up
-- +goose StatementBegin
 ALTER TABLE bookings DROP CONSTRAINT unique_user_workshop;
   CREATE UNIQUE INDEX uniq_user_workshop_confirmed 
     ON bookings(user_id, workshop_id) 
     WHERE status = 'Confirmed';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS uniq_user_workshop_confirmed;
   ALTER TABLE bookings ADD CONSTRAINT unique_user_workshop UNIQUE (user_id, workshop_id);
-- +goose StatementEnd

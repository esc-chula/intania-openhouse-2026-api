-- +goose Up
-- +goose StatementBegin
CREATE TYPE workshop_category AS ENUM ('Department', 'Club');
CREATE TABLE workshops (
    id BIGSERIAL PRIMARY KEY,
    
    name TEXT NOT NULL,
    description TEXT,
    category workshop_category NOT NULL,
    affiliation TEXT NOT NULL,
    
    event_date DATE NOT NULL
        CHECK (event_date IN ('2026-03-28', '2026-03-29')),
    
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    
    location TEXT,
    
    total_seats INTEGER NOT NULL CHECK (total_seats >= 1),
    registered_count INTEGER NOT NULL DEFAULT 0 CHECK (registered_count >= 0),
    
    CHECK (end_time > start_time)
);
CREATE INDEX idx_workshops_event_date ON workshops(event_date);


CREATE TYPE booking_status AS ENUM ('Confirmed', 'Cancelled');
CREATE TABLE bookings (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    workshop_id BIGINT NOT NULL,
    status booking_status NOT NULL DEFAULT 'Confirmed',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_booking_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_booking_workshop
        FOREIGN KEY (workshop_id)
        REFERENCES workshops(id)
        ON DELETE CASCADE,
    
    CONSTRAINT unique_user_workshop
        UNIQUE (user_id, workshop_id)
);
CREATE INDEX idx_bookings_workshop_id ON bookings(workshop_id);
CREATE INDEX idx_bookings_user_id ON bookings(user_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS workshops;
DROP TYPE IF EXISTS booking_status;
DROP TYPE IF EXISTS workshop_category;
-- +goose StatementEnd

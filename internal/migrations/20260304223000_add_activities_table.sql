-- +goose Up
-- +goose StatementBegin
CREATE TABLE activities (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    building_name TEXT,
    floor TEXT,
    room_name TEXT,
    image TEXT,
    
    CHECK (end_time > start_time)
);

CREATE INDEX idx_activities_start_time ON activities(start_time);
CREATE INDEX idx_activities_end_time ON activities(end_time);
CREATE INDEX idx_activities_title ON activities(title);
CREATE INDEX idx_activities_location ON activities(building_name, room_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS activities;
-- +goose StatementEnd

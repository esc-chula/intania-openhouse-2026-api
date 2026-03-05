-- +goose Up
-- +goose StatementBegin

-- Add UNIQUE so that there's no duplicate and also create the unique index for faster lookup
ALTER TABLE workshops ADD COLUMN IF NOT EXISTS check_in_code UUID UNIQUE DEFAULT gen_random_uuid();

CREATE TABLE booths (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    check_in_code UUID UNIQUE DEFAULT gen_random_uuid()
);

CREATE TABLE booth_checkins (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    booth_id BIGINT NOT NULL REFERENCES booths(id) ON DELETE CASCADE,
    checked_in_at TIMESTAMP WITH TIME ZONE NOT NULL,

    UNIQUE (user_id, booth_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS booth_checkins;
DROP TABLE IF EXISTS booths;

ALTER TABLE workshops DROP COLUMN IF EXISTS check_in_code;
-- +goose StatementEnd

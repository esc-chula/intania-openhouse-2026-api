-- +goose Up
-- +goose StatementBegin
CREATE stamp_type AS ENUM ('Department', 'Club', 'Exhibition');
CREATE TABLE stamp_posters (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type stamp_type NOT NULL,
    is_redeemed BOOL NOT NULL DEFAULT FALSE,

    UNIQUE (user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stamp_posters;
-- +goose StatementEnd

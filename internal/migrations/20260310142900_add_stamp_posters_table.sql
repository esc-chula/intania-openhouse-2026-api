-- +goose Up
-- +goose StatementBegin
CREATE TABLE stamp_posters (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stamp_type TEXT NOT NULL,
    is_redeemed BOOL NOT NULL DEFAULT FALSE,

    UNIQUE (user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stamp_posters;
-- +goose StatementEnd

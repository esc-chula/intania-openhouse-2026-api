-- +goose Up
-- +goose StatementBegin
CREATE TYPE stamp_type AS ENUM ('department', 'club', 'exhibition');
CREATE TABLE stamp_posters (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type stamp_type NOT NULL,
    is_redeemed BOOL NOT NULL DEFAULT FALSE,

    UNIQUE (user_id, type)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stamp_posters;
DROP TYPE IF EXISTS stamp_type;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin

ALTER TYPE participant_type_enum RENAME TO participant_type_enum_old;

CREATE TYPE participant_type_enum AS ENUM (
    'student',
    'intania',
    'outside_student',
    'alumni',
    'teacher',
    'other'
);

ALTER TABLE users ALTER COLUMN participant_type TYPE participant_type_enum USING participant_type::text::participant_type_enum;

DROP TYPE participant_type_enum_old;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TYPE participant_type_enum RENAME TO participant_type_enum_old;

CREATE TYPE participant_type_enum AS ENUM (
    'student',
    'intania',
    'other_university_student',
    'teacher',
    'other'
);

ALTER TABLE users ALTER COLUMN participant_type TYPE participant_type_enum USING participant_type::text::participant_type_enum;

DROP TYPE participant_type_enum_old;

-- +goose StatementEnd

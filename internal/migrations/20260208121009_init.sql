-- +goose Up
-- +goose StatementBegin

CREATE TYPE gender_enum AS ENUM (
    'male',              -- ชาย
    'female',            -- หญิง
    'prefer_not_to_say', -- ไม่ต้องการระบุ
    'other'              -- อื่นๆ
);

CREATE TYPE participant_type_enum AS ENUM (
    'student',                  -- นักเรียน/ผู้ที่สนใจศึกษาต่อ
    'intania',                  -- นิสิตปัจจุบัน/นิสิตเก่าวิศวะจุฬาฯ
    'other_university_student', -- นิสิตจากมหาลัยอื่น
    'teacher',                  -- ครู
    'other'                     -- ผู้ปกครอง/บุคคลภายนอก
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    gender gender_enum NOT NULL,
    phone_number TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,

    participant_type participant_type_enum NOT NULL,
    attendance_dates DATE[] DEFAULT '{}',
    interested_activities TEXT[] DEFAULT '{}',
    discovery_channel TEXT[] DEFAULT '{}',
    extra_attributes JSONB DEFAULT '{}'::jsonb,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_set_updated_at_users BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_set_updated_at_users ON users;

DROP FUNCTION IF EXISTS set_updated_at();

DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS participant_type_enum;
DROP TYPE IF EXISTS gender_enum;
-- +goose StatementEnd

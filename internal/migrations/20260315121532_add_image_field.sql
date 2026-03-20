-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'workshops' AND column_name = 'image'
  ) THEN
    ALTER TABLE workshops ADD COLUMN image TEXT;
  END IF;
END$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'workshops' AND column_name = 'image'
  ) THEN
    ALTER TABLE workshops DROP COLUMN image;
  END IF;
END$$;
-- +goose StatementEnd

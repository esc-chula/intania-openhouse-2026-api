-- +goose Up
-- +goose StatementBegin
ALTER TYPE workshop_category RENAME VALUE 'Department' TO 'department';
ALTER TYPE workshop_category RENAME VALUE 'Club' TO 'club';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TYPE workshop_category RENAME VALUE 'department' TO 'Department';
ALTER TYPE workshop_category RENAME VALUE 'club' TO 'Club';
-- +goose StatementEnd

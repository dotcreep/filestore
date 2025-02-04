-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE files ADD COLUMN label_name TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE files DROP COLUMN label_name;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE files ADD COLUMN package_name TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE files DROP COLUMN package_name;
-- +goose StatementEnd

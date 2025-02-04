-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE files ADD COLUMN version_app TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE files DROP COLUMN version_app;
-- +goose StatementEnd

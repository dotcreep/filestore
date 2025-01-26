-- +goose Up
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS files (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				index_app INTEGER NOT NULL,
				upload_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				user_id TEXT NOT NULL,
				hash TEXT NOT NULL,
				filename TEXT NOT NULL, 
				app_name TEXT NOT NULL,
				type_app TEXT NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
-- +goose Down
SELECT 'down SQL query';
DROP TABLE files;
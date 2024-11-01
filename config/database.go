package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase(dbName string) {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		db, err := sql.Open("sqlite3", dbName)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		_, err = db.Exec(`
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
		`)
		if err != nil {
			log.Fatal(err)
		}
	}
}

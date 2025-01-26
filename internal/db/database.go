package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/dotcreep/filestore/internal/utils"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func initDatabase(dbName string) (*sql.DB, error) {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		file, err := os.Create(dbName)
		if err != nil {
			return nil, err
		}
		file.Close()
		log.Printf("database file '%s' created.\n", dbName)
	}
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Println("error opening database:", err)
		return nil, err
	}
	return db, nil
}

func Migration() error {
	cfg, err := utils.OpenYAML()
	if err != nil {
		return err
	}
	goose.SetBaseFS(migrationsFS)
	db, err := initDatabase(cfg.Config.DBName)
	if err != nil {
		return fmt.Errorf("failed initialize database: %v", err)
	}
	defer db.Close()

	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed set dialect: %v", err)
	}
	migrationPath := "migrations"
	log.Println("running migrations")
	err = goose.Up(db, migrationPath)
	if err != nil {
		return fmt.Errorf("failed running migrations database: %v", err)
	}
	log.Println("done running migrations")
	return nil
}

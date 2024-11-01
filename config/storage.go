package config

import (
	"log"
	"os"
)

func InitStorage() {
	var storagePath string
	err := os.Getenv("STORAGE_PATH")
	if err == "" {
		log.Fatal(err)
	}
	storagePath = os.Getenv("STORAGE_PATH")
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		if err := os.MkdirAll(storagePath, 0755); err != nil {
			log.Fatal(err)
		}
	}
}

package utils

import (
	"log"
	"os"
)

func InitStorage() {
	cfg, err := OpenYAML()
	if err != nil {
		log.Println(err)
		return
	}
	storagePath := cfg.Config.StoragePath
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		if err := os.MkdirAll(storagePath, 0755); err != nil {
			log.Println(err)
		}
	}
}

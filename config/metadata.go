package config

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/dotcreep/filestore/env"
)

func GenMetadata(id string, file io.Reader, filename string, index int, dateTime time.Time) (map[string]interface{}, error) {
	hash := sha256.New()
	errEnv := env.Load(".env")
	if errEnv != nil {
		return map[string]interface{}{}, errors.New("failed to load .env file")
	}
	key := os.Getenv("KEY_ACCESS_METADATA")
	combineData := fmt.Sprintf("%s:%s:%s:%s", id, key, dateTime, filename)
	_, err := hash.Write([]byte(combineData))
	if err != nil {
		return map[string]interface{}{}, err
	}
	_, errhash := io.Copy(hash, file)
	if errhash != nil {
		return map[string]interface{}{}, errhash
	}
	hashString := hex.EncodeToString(hash.Sum(nil))
	url := fmt.Sprintf("/%s/%s", id, hashString)
	return map[string]interface{}{
		"index":     index,
		"upload_at": dateTime,
		"filename":  file,
		"hash":      hashString,
		"url":       url,
	}, nil
}

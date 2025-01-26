package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"
)

type Metadata struct {
	Index    int
	UploadAt time.Time
	Filename string
	Hash     string
	Url      string
	AppId    string
}

func GenMetadata(id string, file io.Reader, filename string, index int, packageName string, dateTime time.Time) (*Metadata, error) {
	hash := sha256.New()
	cfg, err := OpenYAML()
	if err != nil {
		return &Metadata{}, errors.New("failed load config file")
	}
	key := cfg.Config.KeyAccessMetadata
	combineData := fmt.Sprintf("%s:%s:%s:%s", id, key, dateTime, filename)
	_, err = hash.Write([]byte(combineData))
	if err != nil {
		return &Metadata{}, err
	}
	_, errhash := io.Copy(hash, file)
	if errhash != nil {
		return &Metadata{}, errhash
	}
	hashString := hex.EncodeToString(hash.Sum(nil))
	url := fmt.Sprintf("/%s/%s", id, hashString)
	return &Metadata{
		Index:    index,
		UploadAt: dateTime,
		Filename: filename,
		Hash:     hashString,
		Url:      url,
		AppId:    packageName,
	}, nil
}

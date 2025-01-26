package utils

import "github.com/gabriel-vasile/mimetype"

func DetectContentType(filePath string) (string, error) {
	mtype, err := mimetype.DetectFile(filePath)
	if err != nil {
		return "", err
	}
	return mtype.String(), nil
}

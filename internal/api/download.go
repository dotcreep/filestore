package api

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dotcreep/filestore/internal/utils"
)

// func findFilenameByFile(userId string, hash string) (string, error) {
// 	userStorage := filepath.Join(storagePath, userId)
// 	if _, err := os.Stat(userStorage); os.IsNotExist(err) {
// 		return "", err
// 	}
// 	filePath := filepath.Join(userStorage, hash)
// 	if _, err := os.Stat(filePath); os.IsNotExist(err) {
// 		return "", err
// 	}
// 	return filePath, nil
// }

func findFilename(userId string, hash string) (string, string, error) {
	dbName := "./file.db"
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return "", "", errors.New("failed to open database")
	}
	defer db.Close()

	// Gabungkan kedua query menjadi satu
	var hashFilename, filename string
	err = db.QueryRow("SELECT app_name, filename FROM files WHERE user_id = ? AND hash = ?", userId, hash).Scan(&hashFilename, &filename)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", errors.New("file not found")
		}
		return "", "", errors.New("failed to query database")
	}

	return hashFilename, filename, nil
}

func Download(w http.ResponseWriter, r *http.Request) {
	Json := utils.Json{}
	cfg, err := utils.OpenYAML()
	if err != nil {
		Json.NewResponse(false, w, nil, "internal server error", http.StatusInternalServerError, err.Error())
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		Json.NewResponse(false, w, nil, "invalid url", http.StatusBadRequest, nil)
		return
	}
	userId := parts[2]
	hash := parts[3]

	storagePath := cfg.Config.StoragePath
	if storagePath == "" {
		Json.NewResponse(false, w, nil, "storage path not configured", http.StatusInternalServerError, nil)
		return
	}

	//_, _ = findFilenameByFile(userId, hash)
	// if err != nil {
	// 	utils.RespondJSON(w, http.StatusNotFound, map[string]interface{}{
	// 		"success": false,
	// 		"message": "File not found",
	// 	})
	// 	return
	// }

	// Get from database
	hashFilename, filename, err := findFilename(userId, hash)
	if err != nil {
		Json.NewResponse(false, w, nil, "file not found", http.StatusNotFound, err.Error())
		return
	}
	// End get from database

	filePath := filepath.Join(storagePath, userId, hashFilename)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		Json.NewResponse(false, w, nil, "file not found", http.StatusNotFound, err.Error())
		return
	}
	fileSize := fileInfo.Size()
	log.Printf("File size: %d bytes", fileSize)

	file, err := os.Open(filePath)
	if err != nil {
		Json.NewResponse(false, w, nil, "file not found", http.StatusNotFound, err.Error())
		return
	}
	defer file.Close()
	ext := filepath.Ext(filename)
	contentType := "application/octet-stream"
	switch ext {
	case ".apk":
		contentType = "application/vnd.android.package-archive"
	case ".aab":
		contentType = "application/octet-stream"
	case ".zip":
		contentType = "application/zip"
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileSize))
	_, err = io.Copy(w, file)
	if err != nil {
		Json.NewResponse(false, w, nil, "failed to download file", http.StatusInternalServerError, err.Error())
		return
	}

	Json.NewResponse(true, w, nil, "File downloaded successfully", http.StatusOK, nil)
}

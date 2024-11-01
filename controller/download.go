package controller

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

	"github.com/dotcreep/filestore/utils"
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
	hashName, err := db.Query("SELECT app_name FROM files WHERE user_id = ? AND hash = ?", userId, hash)
	if err != nil {
		return "", "", errors.New("failed to query database")
	}
	defer hashName.Close()
	var hashFilename string
	for hashName.Next() {
		err := hashName.Scan(&hashFilename)
		if err != nil {
			return "", "", errors.New("failed to scan row")
		}
	}
	frows, err := db.Query("SELECT filename FROM files WHERE user_id = ? AND hash = ?", userId, hash)
	if err != nil {
		return "", "", errors.New("failed to query database")
	}
	defer frows.Close()
	var filename string
	for frows.Next() {
		err := frows.Scan(&filename)
		if err != nil {
			return "", "", errors.New("failed to scan row")
		}
	}
	return hashFilename, filename, nil
}

func Download(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
			"success": false,
			"message": "Method not allowed",
		})
		return
	}
	vars := r.URL.Path
	parts := strings.Split(vars, "/")
	userId := parts[2]
	hash := parts[3]

	var storagePath string
	errStrg := os.Getenv("STORAGE_PATH")
	if errStrg != "" {
		log.Fatal(errStrg)
	}
	storagePath = os.Getenv("STORAGE_PATH")

	//_, _ = findFilenameByFile(userId, hash)
	// if err != nil {
	// 	utils.RespondJSON(w, http.StatusNotFound, map[string]interface{}{
	// 		"success": false,
	// 		"message": "File not found",
	// 	})
	// 	return
	// }

	// Get from database
	hashName, filename, err := findFilename(userId, hash)
	if err != nil {
		utils.RespondJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": "File not found",
		})
		return
	}
	// End get from database

	base := filepath.Join(storagePath, userId)
	//ext := filepath.Ext(filename)
	filePath := filepath.Join(base, hashName)
	file, err := os.Open(filePath)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to open file",
		})
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	_, err = io.Copy(w, file)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to download file",
		})
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "File downloaded successfully",
	})
}

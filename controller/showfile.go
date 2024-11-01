package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dotcreep/filestore/utils"
)

func findFilenameViaStorage(userId string) (string, error) {

	var storagePath string
	err := os.Getenv("STORAGE_PATH")
	if err != "" {
		log.Fatal(err)
	}
	storagePath = os.Getenv("STORAGE_PATH")
	userStorage := filepath.Join(storagePath, userId)
	if _, err := os.Stat(userStorage); os.IsNotExist(err) {
		return "", err
	}
	return userStorage, nil
}

func ShowFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
			"success": false,
			"message": "Method not allowed",
			"data":    map[string]interface{}{},
		})
		return
	}
	vars := r.URL.Path
	parts := strings.Split(vars, "/")
	userId := parts[3]
	dirs, err := findFilenameViaStorage(userId)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			utils.RespondJSON(w, http.StatusNotFound, map[string]interface{}{
				"success": false,
				"message": "User not found",
				"data":    map[string]interface{}{},
			})
			return
		} else {
			utils.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"message": "Internal Server Error",
				"data":    map[string]interface{}{},
			})
			return
		}
	}
	dir, err := os.ReadDir(dirs)
	if err != nil {
		utils.RespondJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": "Failed get directory",
			"data":    map[string]interface{}{},
		})
		return
	}
	dbName := "./file.db"
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to open database",
			"data":    map[string]interface{}{},
		})
		return
	}
	defer db.Close()
	data := make(map[string]interface{})
	data["apk"] = make(map[string]interface{})
	data["aab"] = make(map[string]interface{})
	for _, f := range dir {
		if f.Type().IsRegular() {
			filename := f.Name()
			ext := filepath.Ext(filename)
			if ext == "" {
				utils.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
					"success": false,
					"message": "Failed to get file extension",
					"data":    map[string]interface{}{},
				})
				return
			}
			dbName := "./file.db"
			db, err := sql.Open("sqlite3", dbName)
			if err != nil {
				utils.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
					"success": false,
					"message": "Failed to open database",
					"data":    map[string]interface{}{},
				})
				return
			}
			defer db.Close()
			var fileName string
			err = db.QueryRow("SELECT filename FROM files WHERE user_id = ? AND app_name = ?", userId, filename).Scan(&fileName)
			if err != nil {
				if err == sql.ErrNoRows {
					utils.RespondJSON(w, http.StatusNotFound, map[string]interface{}{
						"success": false,
						"message": "Index not found",
						"data":    map[string]interface{}{},
					})
					return
				} else {
					utils.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
						"success": false,
						"message": "Failed to query filename in database",
						"data":    map[string]interface{}{},
					})
					return
				}
			}
			var hashStr string
			err = db.QueryRow("SELECT hash FROM files WHERE user_id = ? AND app_name = ?", userId, filename).Scan(&hashStr)
			if err != nil {
				if err == sql.ErrNoRows {
					utils.RespondJSON(w, http.StatusNotFound, map[string]interface{}{
						"success": false,
						"message": "Hash not found",
						"data":    map[string]interface{}{},
					})
					return
				} else {
					utils.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
						"success": false,
						"message": "Failed to query hash in database",
						"data":    map[string]interface{}{},
					})
					return
				}
			}
			var index string
			err = db.QueryRow("SELECT index_app FROM files WHERE user_id = ? AND app_name = ?", userId, filename).Scan(&index)
			if err != nil {
				if err == sql.ErrNoRows {
					utils.RespondJSON(w, http.StatusNotFound, map[string]interface{}{
						"success": false,
						"message": "Index not found",
						"data":    map[string]interface{}{},
					})
					return
				} else {
					utils.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
						"success": false,
						"message": "Failed to query index in database",
						"data":    map[string]interface{}{},
					})
					return
				}
			}
			var uploadDate string
			err = db.QueryRow("SELECT upload_at FROM files WHERE user_id = ? AND app_name = ?", userId, filename).Scan(&uploadDate)
			if err != nil {
				if err == sql.ErrNoRows {
					utils.RespondJSON(w, http.StatusNotFound, map[string]interface{}{
						"success": false,
						"message": "Upload date not found",
						"data":    map[string]interface{}{},
					})
					return
				} else {
					utils.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
						"success": false,
						"message": "Failed to query upload_at in database",
						"data":    map[string]interface{}{},
					})
					return
				}
			}
			fileData := map[string]interface{}{
				"index":     index,
				"filename":  fileName,
				"url":       fmt.Sprintf("/%s/%s", userId, hashStr),
				"upload_at": uploadDate,
			}

			if ext == ".apk" {
				data["apk"].(map[string]interface{})[filename] = fileData
			}
			if ext == ".aab" {
				data["aab"].(map[string]interface{})[filename] = fileData
			}

		}
	}
	if len(data) == 0 {
		utils.RespondJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": "User data is empty",
			"data":    map[string]interface{}{},
		})
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Success getting data",
		"data":    data,
	})
}

package controller

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dotcreep/filestore/config"
	"github.com/dotcreep/filestore/utils"

	_ "github.com/mattn/go-sqlite3"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
			"success": false,
			"message": "Method not allowed",
			"url":     "",
		})
		return
	}

	var storagePath string
	errStrg := os.Getenv("STORAGE_PATH")
	if errStrg != "" {
		log.Fatal(errStrg)
	}
	storagePath = os.Getenv("STORAGE_PATH")

	userId := r.URL.Path
	userId = strings.Split(userId, "/")[2]
	if userId == "" {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "User ID is required",
			"url":     "",
		})
		return
	}
	userStorage := filepath.Join(storagePath, userId)
	if _, err := os.Stat(userStorage); os.IsNotExist(err) {
		if err := os.MkdirAll(userStorage, 0755); err != nil {
			utils.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"message": "Failed to create user storage",
				"url":     "",
			})
			return
		}
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Failed to get file",
			"url":     "",
		})
		return
	}
	defer file.Close()

	// check file type
	fileTypeExt := filepath.Ext(header.Filename)
	if header.Filename != "" {
		if fileTypeExt != ".aab" && fileTypeExt != ".apk" {
			utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "Unsupported file type",
				"url":     "",
			})
			return
		}
	} else {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "No file selected",
			"url":     "",
		})
		return
	}
	filename := header.Filename // fmt.Sprintf("%s-%s", userId, header.Filename)
	filePath := filepath.Join(userStorage, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to create file",
			"url":     "",
		})
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to copy file",
			"url":     "",
		})
		return
	}
	// save to database
	dbName := "./file.db"
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to open database",
			"url":     "",
		})
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO files (index_app, upload_at, user_id, hash, filename, app_name, type_app) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to insert file into database",
			"url":     "",
		})
		return
	}
	defer stmt.Close()

	userDir, err := os.ReadDir(userStorage)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to get user directory",
			"url":     "",
		})
		return
	}
	index := 0
	aabFound := 0
	apkFound := 0
	for _, file := range userDir {
		if file.IsDir() {
			continue
		}
		fileInfo, err := file.Info()
		if err != nil {
			continue
		}
		ext := filepath.Ext(fileInfo.Name())
		switch ext {
		case ".aab":
			aabFound++
		case ".apk":
			apkFound++
		}
	}

	// generate metadata
	dateTime := time.Now()
	if fileTypeExt == ".aab" {
		index = aabFound
	}
	if fileTypeExt == ".apk" {
		index = apkFound
	}
	data, err := config.GenMetadata(userId, file, filename, index, dateTime)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to create metadata URL",
			"url":     "",
		})
		return
	}

	hashString, ok := data["hash"].(string)
	if !ok {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Failed get hash data",
			"url":     "",
		})
		return
	}
	uploadDate, ok := data["upload_at"].(time.Time)
	if !ok {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Failed get upload_at data",
			"url":     "",
		})
		return
	}
	shortHash := hashString[:8]
	ext := filepath.Ext(filename)
	hashApp := fmt.Sprintf("%s%s", shortHash, ext)
	errRename := os.Rename(filePath, filepath.Join(userStorage, hashApp))
	if errRename != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to rename file",
			"url":     "",
		})
		return
	}
	_, err = stmt.Exec(index, uploadDate, userId, hashString, filename, hashApp, fileTypeExt)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to insert file into database",
			"url":     "",
		})
		return
	}
	// end of saving to database
	url, ok := data["url"].(string)
	if !ok {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Failed get url data",
			"url":     "",
		})
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "File uploaded successfully",
		"url":     url,
	})
}

package api

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dotcreep/filestore/internal/utils"

	_ "github.com/mattn/go-sqlite3"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	Json := utils.Json{}
	cfg, err := utils.OpenYAML()
	if err != nil {
		Json.NewResponse(false, w, nil, "internal server error", http.StatusInternalServerError, err.Error())
		return
	}

	storagePath := cfg.Config.StoragePath
	if storagePath == "" {
		Json.NewResponse(false, w, nil, "internal server error", http.StatusInternalServerError, "storage path is not set")
		return
	}

	userId := r.URL.Path
	if userIdParts := strings.Split(userId, "/"); len(userIdParts) >= 5 && userIdParts[4] != "" {
		userId = userIdParts[4]
	} else {
		Json.NewResponse(false, w, nil, "user id is required", http.StatusBadRequest, nil)
		return
	}

	packageName := r.FormValue("id")
	if packageName == "" {
		Json.NewResponse(false, w, nil, "app id is required", http.StatusBadRequest, "no package name")
		return
	}

	userStorage := filepath.Join(storagePath, userId)
	if _, err := os.Stat(userStorage); os.IsNotExist(err) {
		if err := os.MkdirAll(userStorage, 0755); err != nil {
			Json.NewResponse(false, w, nil, "failed to create user storage", http.StatusInternalServerError, err.Error())
			return
		}
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		Json.NewResponse(false, w, nil, "failed to get file", http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	fileTypeExt := filepath.Ext(header.Filename)
	if header.Filename == "" {
		Json.NewResponse(false, w, nil, "no file selected", http.StatusBadRequest, nil)
		return
	}

	if fileTypeExt != ".aab" && fileTypeExt != ".apk" {
		Json.NewResponse(false, w, nil, "unsupported file type", http.StatusBadRequest, nil)
		return
	}

	filename := header.Filename
	filePath := filepath.Join(userStorage, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		Json.NewResponse(false, w, nil, "failed to create file", http.StatusInternalServerError, err.Error())
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		Json.NewResponse(false, w, nil, "failed to copy file", http.StatusInternalServerError, err.Error())
		return
	}

	dbName := "./file.db"
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		Json.NewResponse(false, w, nil, "failed to open database", http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO files (index_app, upload_at, user_id, hash, filename, app_name, package_name, type_app) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		Json.NewResponse(false, w, nil, "failed to prepare statement", http.StatusInternalServerError, err.Error())
		return
	}
	defer stmt.Close()

	userDir, err := os.ReadDir(userStorage)
	if err != nil {
		Json.NewResponse(false, w, nil, "failed to get user directory", http.StatusInternalServerError, err.Error())
		return
	}

	aabFound := 0
	apkFound := 0
	for _, fileInfo := range userDir {
		if fileInfo.IsDir() {
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

	index := 0
	switch fileTypeExt {
	case ".aab":
		index = aabFound
	case ".apk":
		index = apkFound
	}

	dateTime := time.Now()
	data, err := utils.GenMetadata(userId, file, filename, index, packageName, dateTime)
	if err != nil {
		Json.NewResponse(false, w, nil, "failed to generate metadata", http.StatusInternalServerError, err.Error())
		return
	}

	hashString := data.Hash
	if hashString == "" {
		Json.NewResponse(false, w, nil, "failed to get hash data", http.StatusBadRequest, "no hash data")
		return
	}

	uploadDate := data.UploadAt
	if uploadDate.IsZero() {
		Json.NewResponse(false, w, nil, "failed to get upload_at data", http.StatusBadRequest, "no upload_at data")
		return
	}

	pkgName := data.AppId
	if pkgName == "" {
		Json.NewResponse(false, w, nil, "failed to get package name", http.StatusBadRequest, "no package name")
		return
	}

	shortHash := hashString[:8]
	ext := filepath.Ext(filename)
	hashApp := fmt.Sprintf("%s%s", shortHash, ext)
	errRename := os.Rename(filePath, filepath.Join(userStorage, hashApp))
	if errRename != nil {
		Json.NewResponse(false, w, nil, "failed to rename file", http.StatusInternalServerError, errRename.Error())
		return
	}
	_, err = stmt.Exec(index, uploadDate, userId, hashString, filename, hashApp, pkgName, fileTypeExt)
	if err != nil {
		Json.NewResponse(false, w, nil, "failed to insert file into database", http.StatusInternalServerError, err.Error())
		return
	}

	url := data.Url
	if url == "" {
		Json.NewResponse(false, w, nil, "failed to get URL data", http.StatusBadRequest, "no URL data")
		return
	}

	Json.NewResponse(true, w, url, "file uploaded successfully", http.StatusOK, nil)
}

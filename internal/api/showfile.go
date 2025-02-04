package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dotcreep/filestore/internal/utils"
)

type Data struct {
	AAB map[string]FileData `json:"aab"`
	APK map[string]FileData `json:"apk"`
}

type FileData struct {
	Index       int    `json:"index"`
	Filename    string `json:"filename"`
	URL         string `json:"url"`
	UploadAt    string `json:"upload_at"`
	LabelName   string `json:"label_name"`
	Version     string `json:"version"`
	PackageName string `json:"package_name"`
}

func findFilenameViaStorage(userId string) (string, error) {
	cfg, err := utils.OpenYAML()
	if err != nil {
		return "", err
	}
	storagePath := cfg.Config.StoragePath
	userStorage := filepath.Join(storagePath, userId)
	if _, err := os.Stat(userStorage); os.IsNotExist(err) {
		return "", err
	}
	return userStorage, nil
}

// @Summary		Show file from user
// @Description	Show all app list of apk and aab
// @Tags			File
// @Accept			json
// @Produce		json
//
// @Param			id	path	string	true	"user id"
//
// @Security		X-API-Key
// @Success		200	{object}	utils.Showfile				"Success"
// @Failure		400	{object}	utils.BadRequest			"Bad request"
// @Failure		500	{object}	utils.InternalServerError	"Internal server error"
// @Router			/api/v1/list/{id} [get]
func ShowFile(w http.ResponseWriter, r *http.Request) {
	Json := utils.Json{}
	cfg, err := utils.OpenYAML()
	if err != nil {
		Json.NewResponse(false, w, nil, "internal server error", http.StatusInternalServerError, err.Error())
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	userId := parts[4]
	dirs, err := findFilenameViaStorage(userId)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			Json.NewResponse(false, w, nil, "user not found", http.StatusNotFound, err.Error())
			return
		} else {
			Json.NewResponse(false, w, nil, "internal server error", http.StatusInternalServerError, err.Error())
			return
		}
	}
	dir, err := os.ReadDir(dirs)
	if err != nil {
		Json.NewResponse(false, w, nil, "failed to get directory", http.StatusInternalServerError, err.Error())
		return
	}

	data := &Data{
		AAB: make(map[string]FileData),
		APK: make(map[string]FileData),
	}
	for _, f := range dir {
		if f.Type().IsRegular() {
			filename := f.Name()
			ext := filepath.Ext(filename)
			if ext == "" {
				Json.NewResponse(false, w, nil, "failed to get file extension", http.StatusInternalServerError, nil)
				return
			}
			dbName := cfg.Config.DBName
			db, err := sql.Open("sqlite3", dbName)
			if err != nil {
				Json.NewResponse(false, w, nil, "failed to open database", http.StatusInternalServerError, err.Error())
				return
			}
			defer db.Close()
			var fileName string
			fmt.Println(userId, filename)
			err = db.QueryRow("SELECT filename FROM files WHERE user_id = ? AND app_name = ?", userId, filename).Scan(&fileName)
			if err != nil {
				if err == sql.ErrNoRows {
					Json.NewResponse(false, w, nil, "index file not found", http.StatusNotFound, err.Error())
					return
				} else {
					Json.NewResponse(false, w, nil, "failed to query index in database", http.StatusInternalServerError, err.Error())
					return
				}
			}
			var hashStr string
			err = db.QueryRow("SELECT hash FROM files WHERE user_id = ? AND app_name = ?", userId, filename).Scan(&hashStr)
			if err != nil {
				if err == sql.ErrNoRows {
					Json.NewResponse(false, w, nil, "hash not found", http.StatusNotFound, err.Error())
					return
				} else {
					Json.NewResponse(false, w, nil, "failed to query hash in database", http.StatusInternalServerError, err.Error())
					return
				}
			}
			var index int
			err = db.QueryRow("SELECT index_app FROM files WHERE user_id = ? AND app_name = ?", userId, filename).Scan(&index)
			if err != nil {
				if err == sql.ErrNoRows {
					Json.NewResponse(false, w, nil, "index not found", http.StatusNotFound, err.Error())
					return
				} else {
					Json.NewResponse(false, w, nil, "failed to query index in database", http.StatusInternalServerError, err.Error())
					return
				}
			}
			var labelName sql.NullString
			err = db.QueryRow("SELECT label_name FROM files WHERE user_id = ? AND app_name = ?", userId, filename).Scan(&labelName)
			if err != nil {
				if err == sql.ErrNoRows {
					Json.NewResponse(false, w, nil, "label not found", http.StatusNotFound, err.Error())
					return
				} else {
					Json.NewResponse(false, w, nil, "failed to query label in database", http.StatusInternalServerError, err.Error())
					return
				}
			}
			var versionApp sql.NullString
			err = db.QueryRow("SELECT version_app FROM files WHERE user_id = ? AND app_name = ?", userId, filename).Scan(&versionApp)
			if err != nil {
				if err == sql.ErrNoRows {
					Json.NewResponse(false, w, nil, "version not found", http.StatusNotFound, err.Error())
					return
				} else {
					Json.NewResponse(false, w, nil, "failed to query version in database", http.StatusInternalServerError, err.Error())
					return
				}
			}
			var uploadDate string
			err = db.QueryRow("SELECT upload_at FROM files WHERE user_id = ? AND app_name = ?", userId, filename).Scan(&uploadDate)
			if err != nil {
				if err == sql.ErrNoRows {
					Json.NewResponse(false, w, nil, "upload date not found", http.StatusNotFound, err.Error())
					return
				} else {
					Json.NewResponse(false, w, nil, "failed to query upload_at in database", http.StatusInternalServerError, err.Error())
					return
				}
			}
			var pkgName string
			err = db.QueryRow("SELECT package_name FROM files WHERE user_id = ? AND app_name = ?", userId, filename).Scan(&pkgName)
			if err != nil {
				if err == sql.ErrNoRows {
					Json.NewResponse(false, w, nil, "package not found", http.StatusNotFound, err.Error())
					return
				} else {
					Json.NewResponse(false, w, nil, "failed to query package in database", http.StatusInternalServerError, err.Error())
					return
				}
			}
			fileData := &FileData{
				Index:       index,
				Filename:    fileName,
				LabelName:   labelName.String,
				PackageName: pkgName,
				UploadAt:    uploadDate,
				URL:         fmt.Sprintf("/%s/%s", userId, hashStr),
				Version:     versionApp.String,
			}

			if ext == ".apk" {
				data.APK[filename] = *fileData
			}
			if ext == ".aab" {
				data.AAB[filename] = *fileData
			}
		}
	}
	if len(data.AAB) == 0 && len(data.APK) == 0 {
		Json.NewResponse(false, w, nil, "user data is empty", http.StatusOK, nil)
		return
	}
	Json.NewResponse(false, w, data, "success getting data", http.StatusOK, nil)
}

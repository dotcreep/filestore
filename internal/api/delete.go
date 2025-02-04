package api

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dotcreep/filestore/internal/utils"
)

//	@Summary		Delete file apk or aab
//	@Description	Delete file based on user id or username using hash
//	@Tags			File
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string	true	"user id"
//	@Param			hash	path	string	true	"hash"
//	@Security		X-API-Key
//	@Success		200	{object}	utils.Success				"Success"
//	@Failure		400	{object}	utils.BadRequest			"Bad request"
//	@Failure		500	{object}	utils.InternalServerError	"Internal server error"
//	@Router			/api/v1/{id}/{hash} [delete]
func Delete(w http.ResponseWriter, r *http.Request) {
	Json := utils.Json{}
	cfg, err := utils.OpenYAML()
	if err != nil {
		Json.NewResponse(false, w, nil, "internal server error", http.StatusInternalServerError, err.Error())
		return
	}
	storage := cfg.Config.StoragePath
	if storage == "" {
		Json.NewResponse(false, w, nil, "internal server error", http.StatusInternalServerError, "storage path is not set")
		return
	}
	userId := r.URL.Path
	if userIdParts := strings.Split(userId, "/"); len(userIdParts) >= 4 && userIdParts[3] != "" {
		userId = userIdParts[3]
	} else {
		Json.NewResponse(false, w, nil, "user id is required", http.StatusBadRequest, nil)
		return
	}

	hash := r.URL.Path
	if hashParts := strings.Split(hash, "/"); len(hashParts) >= 5 && hashParts[4] != "" {
		hash = hashParts[4]
	} else {
		Json.NewResponse(false, w, nil, "hash is required", http.StatusBadRequest, nil)
		return
	}

	db, err := sql.Open("sqlite3", cfg.Config.DBName)
	if err != nil {
		Json.NewResponse(false, w, nil, "internal server error", http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Close()

	var appName string
	err = db.QueryRow("SELECT app_name FROM files WHERE user_id = ? AND hash = ?", userId, hash).Scan(&appName)
	if err != nil {
		Json.NewResponse(false, w, nil, "file not found", http.StatusNotFound, err.Error())
		return
	}

	_, err = db.Exec("DELETE FROM files WHERE user_id = ? AND hash = ?", userId, hash)
	if err != nil {
		Json.NewResponse(false, w, nil, "failed to delete file", http.StatusInternalServerError, err.Error())
		return
	}

	userStorage := filepath.Join(storage, userId)
	if _, err := os.Stat(userStorage); os.IsNotExist(err) {
		Json.NewResponse(false, w, nil, "user not found", http.StatusNotFound, err.Error())
		return
	}
	filePath := filepath.Join(userStorage, appName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		Json.NewResponse(false, w, nil, "file not found", http.StatusNotFound, err.Error())
		return
	}
	err = os.Remove(filePath)
	if err != nil {
		Json.NewResponse(false, w, nil, "failed to delete file", http.StatusInternalServerError, err.Error())
		return
	}
	Json.NewResponse(true, w, nil, "file deleted", http.StatusOK, nil)
}

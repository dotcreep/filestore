package api

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dotcreep/filestore/internal/utils"
)

func DeleteByUsername(w http.ResponseWriter, r *http.Request) {
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

	db, err := sql.Open("sqlite3", cfg.Config.DBName)
	if err != nil {
		Json.NewResponse(false, w, nil, "internal server error", http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM files WHERE user_id = ?", userId)
	if err != nil {
		Json.NewResponse(false, w, nil, "failed to delete data", http.StatusInternalServerError, err.Error())
		return
	}

	userStorage := filepath.Join(storage, userId)
	if _, err := os.Stat(userStorage); os.IsNotExist(err) {
		Json.NewResponse(false, w, nil, "user not found", http.StatusNotFound, err.Error())
		return
	}
	err = os.RemoveAll(userStorage)
	if err != nil {
		Json.NewResponse(false, w, nil, "failed to delete file", http.StatusInternalServerError, err.Error())
		return
	}
	Json.NewResponse(true, w, nil, "user deleted", http.StatusOK, nil)
}

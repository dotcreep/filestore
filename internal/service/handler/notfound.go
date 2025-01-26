package handler

import (
	"net/http"

	"github.com/dotcreep/filestore/internal/utils"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	Json := utils.Json{}
	Json.NewResponse(false, w, nil, "404 not found", http.StatusNotFound, nil)
}

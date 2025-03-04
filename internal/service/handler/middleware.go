package handler

import (
	"log"
	"net/http"

	"github.com/dotcreep/filestore/internal/utils"
)

func Middleware(next http.Handler) http.Handler {
	Json := utils.Json{}
	cfg, err := utils.OpenYAML()
	if err != nil {
		log.Println(err)
		return nil
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := cfg.Config.SecretAccess
		if secret == "" {
			Json.NewResponse(false, w, nil, "Unauthorized", http.StatusInternalServerError, "failed to get secret from config")
			return
		}
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != secret {
			Json.NewResponse(false, w, nil, "Unauthorized", http.StatusUnauthorized, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	Json := utils.Json{}
	cfg, err := utils.OpenYAML()
	if err != nil {
		log.Println(err)
		return nil
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := cfg.Config.SuperUserKey
		if secret == "" {
			Json.NewResponse(false, w, nil, "Unauthorized", http.StatusInternalServerError, "failed to get secret from config")
			return
		}
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != secret {
			Json.NewResponse(false, w, nil, "Unauthorized", http.StatusUnauthorized, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

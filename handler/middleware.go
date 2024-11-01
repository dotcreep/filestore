package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/dotcreep/filestore/env"
	"github.com/dotcreep/filestore/utils"
)

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		env.Load(".env")
		secret := os.Getenv("SECRET_ACCESS")
		if secret == "" {
			log.Fatal("Failed to get secret from .env file")
		}
		apiKey := r.Header.Get("X-Auth-Key")
		if apiKey != secret {
			utils.RespondJSON(w, http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"message": "Unauthorized",
			})
			return
		}
		next.ServeHTTP(w, r)
	}
}

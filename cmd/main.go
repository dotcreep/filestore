package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/dotcreep/filestore/docs"
	"github.com/dotcreep/filestore/internal/api"
	database "github.com/dotcreep/filestore/internal/db"
	"github.com/dotcreep/filestore/internal/service/handler"
	"github.com/dotcreep/filestore/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/mattn/go-sqlite3"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title						Filestore
// @version					1.0
// @description				Documentation for Filestore
// @license.name				MIT
// @license.url				https://opensource.org/licenses/MIT
// @BasePath					/
// @SecurityDefinitions.apikey	X-API-Key
// @name						X-API-Key
// @in							header
// @description				Input your token authorized
func main() {
	cfg, err := utils.OpenYAML()
	if err != nil {
		log.Println("failed load config file")
		panic(err)
	}
	port := cfg.Config.Port
	if port == 0 {
		port = 9091
		fmt.Printf("Warning: PORT is not set in config file, make default %d\n", port)
	}
	if cfg.Config.DBName == "" {
		cfg.Config.DBName = "file.db"
		log.Printf("Warning: DBName is not set in config file, make default %s\n", cfg.Config.DBName)
	}

	err = database.Migration()
	if err != nil {
		log.Println(err)
	}
	utils.InitStorage()

	r := chi.NewRouter()
	cors := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-API-Key"},
		AllowCredentials: true,
	})
	r.Use(cors)
	r.Route("/", func(r chi.Router) {
		r.Get("/getfile/{userId}/{hash}", api.Download)
	})
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(handler.Middleware)
		r.Get("/list/{userId}", api.ShowFile)
		r.Post("/upload/{userId}", api.Upload)
	})
	r.Route("/api/v1/admin", func(r chi.Router) {
		r.Use(handler.AdminMiddleware)
		r.Delete("/{userId}/{hash}", api.Delete)
		r.Delete("/{userId}", api.DeleteByUsername)
	})
	r.Get("/docs/*", httpSwagger.WrapHandler)
	r.NotFound(handler.NotFoundHandler)
	fmt.Printf("Server running on port %d\n", port)
	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}

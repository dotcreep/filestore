package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dotcreep/filestore/config"
	"github.com/dotcreep/filestore/controller"
	"github.com/dotcreep/filestore/env"
	"github.com/dotcreep/filestore/handler"
)

func main() {
	m := handler.Middleware
	err := env.Load(".env")
	if err != nil {
		log.Fatal("Failed to load .env file")
	}
	var middleware bool
	middle := os.Getenv("MIDDLEWARE")
	if middle == "" {
		middleware = true
	}
	middleware, _ = strconv.ParseBool(middle)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
		fmt.Printf("PORT is not set in .env file, make default %s\n", port)
	}
	config.InitDatabase("./file.db")
	config.InitStorage()
	mux := http.NewServeMux()
	if middleware {
		mux.HandleFunc("POST /upload/{userId}", m(controller.Upload))
		mux.HandleFunc("GET /getfile/{userId}/{hash}", m(controller.Download))
		mux.HandleFunc("GET /app/list/{userId}", m(controller.ShowFile))
		mux.HandleFunc("/", m(handler.NotFoundHandler))
	} else {
		mux.HandleFunc("POST /upload/{userId}", controller.Upload)
		mux.HandleFunc("GET /getfile/{userId}/{hash}", controller.Download)
		mux.HandleFunc("GET /app/list/{userId}", controller.ShowFile)
		mux.HandleFunc("/", handler.NotFoundHandler)
	}
	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}

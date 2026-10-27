package main

import (
	"fmt"
	"log"
	"net/http"
	"video-catalog/internal/handler"
	"video-catalog/internal/repository"
	m "video-catalog/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load Env", err)
	}

	repository.Connect_DB()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("video catalog is working fine"))
	})
	r.Get("/videos", handler.ListVideo)

	r.Group(func(r chi.Router) {
		r.Use(m.RequestAuth)
		r.Post("/videos/upload", handler.UploadVideoFile)
	})
	fmt.Println("Starting Video Catalog service on :8081...")
	http.ListenAndServe(":8081", r)
}

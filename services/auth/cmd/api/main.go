package main

import (
	"auth-service/internal/handler"
	"auth-service/internal/repository"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	m "auth-service/internal/middleware"
)

func main() {

	//connect to the database
	repository.Connect_DB() // this grants us the excess to the DB that exists in folder repo db.go
	//get a chi router
	r := chi.NewRouter()
	//give it a logger
	r.Use(middleware.Logger)
	//git it a middleware recoverer
	r.Use(middleware.Recoverer)
	//just to confirm that middle is working
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("AUTH Service is running"))
	})
	//post request
	r.Post("/signup", handler.Signup)
	r.Post("/login", handler.Login)

	//middleware protected routes
	r.Group(func(r chi.Router) {
		r.Use(m.RequestAuth)

		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			userId := r.Context().Value("user_id")
			w.WriteHeader(http.StatusOK)
			w.Write(fmt.Appendf(nil, "Welcome to your protected profile! Your User ID is: %v", userId))
		})
	})
	//start the actual server for service
	slog.Info("Starting auth-server at port 8080 ...")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

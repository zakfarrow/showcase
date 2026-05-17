package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"showcase/internal/auth"
	"showcase/internal/database"
	"showcase/internal/handler"
)

func main() {
	_ = godotenv.Load()

	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	auth.InitStore()
	auth.InitOAuth()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Pages
	r.Get("/", handler.Home)
	r.Get("/projects/{slug}", handler.ProjectDetail)

	// Auth
	r.Get("/auth/login", auth.LoginHandler)
	r.Get("/auth/callback", auth.CallbackHandler)
	r.Get("/auth/logout", auth.LogoutHandler)

	// Admin (protected)
	r.Route("/admin", func(r chi.Router) {
		r.Use(auth.RequireAuth)
		r.Get("/", handler.AdminDashboard)
		r.Get("/projects/new", handler.AdminNewProject)
		r.Post("/projects", handler.AdminCreateProject)
		r.Get("/projects/{id}/edit", handler.AdminEditProject)
		r.Post("/projects/{id}", handler.AdminUpdateProject)
		r.Post("/projects/{id}/delete", handler.AdminDeleteProject)
	})

	// API
	r.Route("/api", func(r chi.Router) {
		r.Get("/greeting", handler.Greeting)
		r.Get("/projects", handler.APIListProjects)
		r.Get("/projects/{slug}", handler.APIGetProject)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}

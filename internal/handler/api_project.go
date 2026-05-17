package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"showcase/internal/repository"
)

func APIListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := repository.GetAllProjects(r.Context())
	if err != nil {
		log.Printf("Error fetching projects: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func APIGetProject(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "Project slug required", http.StatusBadRequest)
		return
	}

	project, err := repository.GetProjectBySlug(r.Context(), slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Project not found", http.StatusNotFound)
			return
		}
		log.Printf("Error fetching project: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

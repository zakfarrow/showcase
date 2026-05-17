package handler

import (
	"log"
	"net/http"

	"showcase/internal/model"
	"showcase/internal/repository"
	"showcase/templates/pages"
)

func Home(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	tech := r.URL.Query().Get("tech")

	var projects []model.Project
	var err error

	if status != "" || tech != "" {
		projects, err = repository.GetFilteredProjects(r.Context(), status, tech)
	} else {
		projects, err = repository.GetAllProjects(r.Context())
	}
	if err != nil {
		log.Printf("Error fetching projects: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	statuses, err := repository.GetAllStatuses(r.Context())
	if err != nil {
		log.Printf("Error fetching statuses: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	technologies, err := repository.GetAllTechnologies(r.Context())
	if err != nil {
		log.Printf("Error fetching technologies: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	pages.Home(projects, statuses, technologies, status, tech).Render(r.Context(), w)
}

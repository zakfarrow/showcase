package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"showcase/internal/repository"
	"showcase/templates/pages"
)

func ProjectDetail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	project, err := repository.GetProjectBySlug(r.Context(), slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		log.Printf("Error fetching project: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	pages.ProjectDetail(project).Render(r.Context(), w)
}

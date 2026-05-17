package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"showcase/internal/auth"
	"showcase/internal/model"
	"showcase/internal/repository"
	"showcase/templates/pages/admin"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	projects, err := repository.GetAllProjects(r.Context())
	if err != nil {
		log.Printf("Error fetching projects: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user := auth.UserFromContext(r.Context())
	admin.Dashboard(projects, user).Render(r.Context(), w)
}

func AdminNewProject(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	admin.ProjectForm(nil, user).Render(r.Context(), w)
}

func AdminCreateProject(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	project := projectFromForm(r)
	if err := repository.CreateProject(r.Context(), project); err != nil {
		log.Printf("Error creating project: %v", err)
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func AdminEditProject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	project, err := repository.GetProjectByID(r.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		log.Printf("Error fetching project: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user := auth.UserFromContext(r.Context())
	admin.ProjectForm(project, user).Render(r.Context(), w)
}

func AdminUpdateProject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	project := projectFromForm(r)
	project.ID = id

	if err := repository.UpdateProject(r.Context(), project); err != nil {
		log.Printf("Error updating project: %v", err)
		http.Error(w, "Failed to update project", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func AdminDeleteProject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if err := repository.DeleteProject(r.Context(), id); err != nil {
		log.Printf("Error deleting project: %v", err)
		http.Error(w, "Failed to delete project", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func projectFromForm(r *http.Request) *model.Project {
	sortOrder, _ := strconv.Atoi(r.FormValue("sort_order"))

	p := &model.Project{
		Slug:        r.FormValue("slug"),
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Content:     r.FormValue("content"),
		TechStack:   parseCommaSeparated(r.FormValue("tech_stack")),
		Status:      r.FormValue("status"),
		Features:    parseCommaSeparated(r.FormValue("features")),
		Challenges:  r.FormValue("challenges"),
		Learnings:   r.FormValue("learnings"),
		FuturePlans: r.FormValue("future_plans"),
		Featured:    r.FormValue("featured") == "on",
		SortOrder:   sortOrder,
	}

	if url := r.FormValue("github_url"); url != "" {
		p.GitHubURL = &url
	}
	if url := r.FormValue("live_url"); url != "" {
		p.LiveURL = &url
	}
	if url := r.FormValue("image_url"); url != "" {
		p.ImageURL = &url
	}

	return p
}

func parseCommaSeparated(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

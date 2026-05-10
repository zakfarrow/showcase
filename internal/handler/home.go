package handler

import (
	"net/http"

	"showcase/templates/pages"
)

func Home(w http.ResponseWriter, r *http.Request) {
	pages.Home().Render(r.Context(), w)
}

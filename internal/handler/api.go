package handler

import (
	"net/http"
	"time"

	"showcase/templates/components"
)

func Greeting(w http.ResponseWriter, r *http.Request) {
	components.Greeting(time.Now().Format("15:04:05")).Render(r.Context(), w)
}

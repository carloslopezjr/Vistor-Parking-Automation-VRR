package handlers

import (
	"html/template"
	"net/http"

	"vistor-parking-automation-vrr/internal/models"
)

// HomeHandler renders the home page with vehicle selector.
type HomeHandler struct {
	Vehicles []models.Vehicle
	Tpl      *template.Template
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type homeView struct {
		Vehicles []models.Vehicle
	}

	data := homeView{
		Vehicles: h.Vehicles,
	}

	if err := h.Tpl.ExecuteTemplate(w, "home.html", data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

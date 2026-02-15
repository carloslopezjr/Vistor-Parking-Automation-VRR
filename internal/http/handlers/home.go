package handlers

import (
	"html/template"
	"net/http"

	"vistor-parking-automation-vrr/internal/store"
)

// HomeHandler renders the home page with profile buttons and last status.
type HomeHandler struct {
	Profiles *store.ProfileStore
	Logs     *store.LogStore
	Tpl      *template.Template
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	profiles, err := h.Profiles.List(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ids := make([]int64, 0, len(profiles))
	for _, p := range profiles {
		ids = append(ids, p.ID)
	}
	latestLogs, err := h.Logs.LatestByProfile(ctx, ids)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	type profileView struct {
		Profile interface{}
		LastLog interface{}
	}

	var data struct {
		Profiles []profileView
	}

	for _, p := range profiles {
		pv := profileView{Profile: p}
		if l, ok := latestLogs[p.ID]; ok {
			pv.LastLog = l
		}
		data.Profiles = append(data.Profiles, pv)
	}

	if err := h.Tpl.ExecuteTemplate(w, "home.html", data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

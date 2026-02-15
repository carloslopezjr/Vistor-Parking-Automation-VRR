package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"html/template"

	"vistor-parking-automation-vrr/internal/models"
	"vistor-parking-automation-vrr/internal/store"
)

// ProfilesHandler handles profile CRUD and list views.
type ProfilesHandler struct {
	Profiles *store.ProfileStore
	Tpl      *template.Template
}

func (h *ProfilesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/profiles")
	if path == "" || path == "/" {
		// /profiles
		if r.Method == http.MethodGet {
			h.list(w, r)
			return
		}
		if r.Method == http.MethodPost {
			h.create(w, r)
			return
		}
	} else if path == "/new" && r.Method == http.MethodGet {
		// /profiles/new
		h.newForm(w, r)
		return
	} else {
		// /profiles/{id}[...]
		trimmed := strings.TrimPrefix(path, "/")
		parts := strings.Split(trimmed, "/")
		if len(parts) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		id, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(parts) == 1 {
			// /profiles/{id}
			if r.Method == http.MethodPost {
				h.update(w, r, id)
				return
			}
		} else if len(parts) == 2 {
			switch parts[1] {
			case "edit":
				if r.Method == http.MethodGet {
					h.editForm(w, r, id)
					return
				}
			case "delete":
				if r.Method == http.MethodPost {
					h.delete(w, r, id)
					return
				}
			case "register":
				if r.Method == http.MethodPost {
					// Registration handler will be wired to automation later.
					// For now, just redirect back to home.
					http.Redirect(w, r, "/", http.StatusSeeOther)
					return
				}
			}
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func (h *ProfilesHandler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	profiles, err := h.Profiles.List(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := struct {
		Profiles []models.Profile
	}{Profiles: profiles}

	if err := h.Tpl.ExecuteTemplate(w, "profiles_list.html", data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *ProfilesHandler) newForm(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Profile *models.Profile
		Action  string
	}{
		Profile: nil,
		Action:  "/profiles",
	}
	if err := h.Tpl.ExecuteTemplate(w, "profile_form.html", data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *ProfilesHandler) editForm(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()
	p, err := h.Profiles.Get(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data := struct {
		Profile *models.Profile
		Action  string
	}{
		Profile: p,
		Action:  "/profiles/" + strconv.FormatInt(id, 10),
	}
	if err := h.Tpl.ExecuteTemplate(w, "profile_form.html", data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *ProfilesHandler) create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p := &models.Profile{
		ApartmentName: r.FormValue("apartment_name"),
		LicensePlate:  strings.ToUpper(strings.TrimSpace(r.FormValue("license_plate"))),
		VehicleMake:   strings.TrimSpace(r.FormValue("vehicle_make")),
		VehicleModel:  strings.TrimSpace(r.FormValue("vehicle_model")),
		ResidentName:  strings.TrimSpace(r.FormValue("resident_name")),
		UnitNumber:    strings.TrimSpace(r.FormValue("unit_number")),
		VisitorName:   strings.TrimSpace(r.FormValue("visitor_name")),
		ResidentEmail: strings.TrimSpace(r.FormValue("resident_email")),
		ProfileName:   strings.TrimSpace(r.FormValue("profile_name")),
	}

	// Simple validation: ensure required fields are present.
	if p.ProfileName == "" || p.ApartmentName == "" || p.ResidentName == "" ||
		p.UnitNumber == "" || p.VisitorName == "" || p.LicensePlate == "" ||
		p.VehicleMake == "" || p.VehicleModel == "" || p.ResidentEmail == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if _, err := h.Profiles.Create(ctx, p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profiles", http.StatusSeeOther)
}

func (h *ProfilesHandler) update(w http.ResponseWriter, r *http.Request, id int64) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p := &models.Profile{
		ID:            id,
		ApartmentName: r.FormValue("apartment_name"),
		LicensePlate:  strings.ToUpper(strings.TrimSpace(r.FormValue("license_plate"))),
		VehicleMake:   strings.TrimSpace(r.FormValue("vehicle_make")),
		VehicleModel:  strings.TrimSpace(r.FormValue("vehicle_model")),
		ResidentName:  strings.TrimSpace(r.FormValue("resident_name")),
		UnitNumber:    strings.TrimSpace(r.FormValue("unit_number")),
		VisitorName:   strings.TrimSpace(r.FormValue("visitor_name")),
		ResidentEmail: strings.TrimSpace(r.FormValue("resident_email")),
		ProfileName:   strings.TrimSpace(r.FormValue("profile_name")),
	}

	ctx := r.Context()
	if err := h.Profiles.Update(ctx, p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profiles", http.StatusSeeOther)
}

func (h *ProfilesHandler) delete(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()
	if err := h.Profiles.Delete(ctx, id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/profiles", http.StatusSeeOther)
}

// ensure we reference time to avoid unused import if not yet expanded.
var _ = time.Now

package handlers

import (
	"net/http"
	"strings"

	"vistor-parking-automation-vrr/internal/automation"
	"vistor-parking-automation-vrr/internal/jobs"
	"vistor-parking-automation-vrr/internal/mailer"
	"vistor-parking-automation-vrr/internal/models"
	"vistor-parking-automation-vrr/internal/store"
	"vistor-parking-automation-vrr/internal/tokens"
)

// TokenHandler processes one-click reminder links.
type TokenHandler struct {
	Tokens    tokens.Service
	Profiles  *store.ProfileStore
	Jobs      jobs.Service
	Mailer    mailer.Service
	BaseURL   string
	Automator automation.Service
}

func (h *TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Expect URL like /r/{token}
	token := strings.TrimPrefix(r.URL.Path, "/r/")
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	profileID, err := h.Tokens.ValidateAndConsume(ctx, token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Ensure the profile still exists before attempting automation.
	if _, err := h.Profiles.Get(ctx, profileID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := h.Automator.RegisterVisitor(ctx, profileID, models.TriggerTokenClick)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = result // will be expanded when automation.Service is fully defined

	// TODO: update last_registration_at, log registration, send success/failure email,
	// and schedule new reminder jobs. This will be implemented after automation
	// and logging wiring are complete.

	w.WriteHeader(http.StatusOK)
}

package handlers

import (
	"encoding/json"
	"net/http"

	"vistor-parking-automation-vrr/internal/automation"
	"vistor-parking-automation-vrr/internal/models"
)

// RegisterHandler handles POST requests to /register and triggers automation.
type RegisterHandler struct {
	Vehicles  []models.Vehicle
	Automator automation.Service
}

// RegisterRequest is the JSON request body for registration.
type RegisterRequest struct {
	VehicleID int `json:"vehicleId"`
}

// RegisterResponse is the JSON response from registration.
type RegisterResponse struct {
	Success bool            `json:"success"`
	Error   string          `json:"error,omitempty"`
	Message string          `json:"message,omitempty"`
	Logs    []string        `json:"logs,omitempty"`
	Vehicle *models.Vehicle `json:"vehicle,omitempty"`
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(RegisterResponse{
			Success: false,
			Error:   "method not allowed",
		})
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RegisterResponse{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	// Validate vehicle ID
	if req.VehicleID < 0 || req.VehicleID >= len(h.Vehicles) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RegisterResponse{
			Success: false,
			Error:   "invalid vehicle ID",
		})
		return
	}

	vehicle := h.Vehicles[req.VehicleID]

	// Call automation service
	result, err := h.Automator.RegisterVisitor(r.Context(), vehicle)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := RegisterResponse{
			Success: false,
			Error:   err.Error(),
			Logs:    result.Logs,
			Vehicle: &vehicle,
		}
		if result != nil {
			resp.Message = result.Error
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	if !result.Success {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RegisterResponse{
			Success: false,
			Error:   result.Error,
			Message: result.Error,
			Logs:    result.Logs,
			Vehicle: &vehicle,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RegisterResponse{
		Success: true,
		Message: "Registration completed successfully",
		Logs:    result.Logs,
		Vehicle: &vehicle,
	})
}

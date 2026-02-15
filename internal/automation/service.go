package automation

import (
	"context"

	"vistor-parking-automation-vrr/internal/models"
)

// Result captures the structured outcome of an automation attempt.
type Result struct {
	Success bool
	Error   string
	Logs    []string
}

// Service defines the automation capabilities of the system.
type Service interface {
	RegisterVisitor(ctx context.Context, profileID int64, trigger models.RegistrationTrigger) (Result, error)
}

// noopService is a placeholder implementation used during early wiring.
type noopService struct{}

// NewNoopService returns a Service that does nothing and always reports
// failure. This is used as a stand-in until the real Playwright-backed
// implementation is provided.
func NewNoopService() Service {
	return &noopService{}
}

func (n *noopService) RegisterVisitor(ctx context.Context, profileID int64, trigger models.RegistrationTrigger) (Result, error) {
	return Result{Success: false, Error: "automation not implemented"}, nil
}

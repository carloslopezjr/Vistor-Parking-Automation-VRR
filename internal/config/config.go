package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"vistor-parking-automation-vrr/internal/automation"
)

// Config holds all environment-driven configuration for the application.
type Config struct {
	Port                      int
	BaseURL                   string
	BasicAuthUser             string // optional
	BasicAuthPass             string // optional
	VehiclesCSVPath           string
	PlaywrightTargetURL       string
	PlaywrightHeadless        bool
	PlaywrightTimeout         time.Duration
	PlaywrightWaitAfterSubmit time.Duration
	PlaywrightSelectors       automation.SiteSelectors
}

// Load reads configuration from environment variables and applies defaults
// where appropriate. It returns an error if mandatory values are missing or
// invalid.
func Load() (*Config, error) {
	c := &Config{}

	// PORT
	portStr := getenvDefault("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 {
		return nil, fmt.Errorf("invalid PORT: %q", portStr)
	}
	c.Port = port

	c.BaseURL = os.Getenv("BASE_URL")
	if c.BaseURL == "" {
		return nil, fmt.Errorf("BASE_URL is required")
	}

	c.VehiclesCSVPath = getenvDefault("VEHICLES_CSV_PATH", "./vehicles.csv")

	c.PlaywrightTargetURL = os.Getenv("PLAYWRIGHT_TARGET_URL")
	if c.PlaywrightTargetURL == "" {
		return nil, fmt.Errorf("PLAYWRIGHT_TARGET_URL is required")
	}

	headlessStr := getenvDefault("PLAYWRIGHT_HEADLESS", "true")
	c.PlaywrightHeadless = headlessStr != "false" && headlessStr != "0"

	// Parse playwright timeout
	timeoutStr := getenvDefault("PLAYWRIGHT_TIMEOUT", "120s")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PLAYWRIGHT_TIMEOUT: %q", timeoutStr)
	}
	c.PlaywrightTimeout = timeout

	// Parse wait after submit
	waitStr := getenvDefault("PLAYWRIGHT_WAIT_AFTER_SUBMIT", "5s")
	wait, err := time.ParseDuration(waitStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PLAYWRIGHT_WAIT_AFTER_SUBMIT: %q", waitStr)
	}
	c.PlaywrightWaitAfterSubmit = wait

	// Parse selectors from environment
	c.PlaywrightSelectors = automation.SiteSelectors{
		ApartmentName:   os.Getenv("PLAYWRIGHT_SELECTOR_APARTMENT_NAME"),
		LicensePlate:    os.Getenv("PLAYWRIGHT_SELECTOR_LICENSE_PLATE"),
		VehicleMake:     os.Getenv("PLAYWRIGHT_SELECTOR_VEHICLE_MAKE"),
		VehicleModel:    os.Getenv("PLAYWRIGHT_SELECTOR_VEHICLE_MODEL"),
		ResidentName:    os.Getenv("PLAYWRIGHT_SELECTOR_RESIDENT_NAME"),
		UnitNumber:      os.Getenv("PLAYWRIGHT_SELECTOR_UNIT_NUMBER"),
		VisitorName:     os.Getenv("PLAYWRIGHT_SELECTOR_VISITOR_NAME"),
		ResidentEmail:   os.Getenv("PLAYWRIGHT_SELECTOR_RESIDENT_EMAIL"),
		SubmitButton:    os.Getenv("PLAYWRIGHT_SELECTOR_SUBMIT_BUTTON"),
		SuccessSelector: os.Getenv("PLAYWRIGHT_SELECTOR_SUCCESS"),
	}

	c.BasicAuthUser = os.Getenv("BASIC_AUTH_USER")
	c.BasicAuthPass = os.Getenv("BASIC_AUTH_PASS")

	return c, nil
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

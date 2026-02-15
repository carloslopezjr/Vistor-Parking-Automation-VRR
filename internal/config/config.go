package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all environment-driven configuration for the application.
type Config struct {
	Port               int
	DatabasePath       string
	BaseURL            string
	SMTPHost           string
	SMTPPort           int
	SMTPUser           string
	SMTPPass           string
	SMTPFrom           string
	SiteTargetURL      string
	PlaywrightHeadless bool
	BasicAuthUser      string // optional
	BasicAuthPass      string // optional
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

	c.DatabasePath = os.Getenv("DATABASE_PATH")
	if c.DatabasePath == "" {
		return nil, fmt.Errorf("DATABASE_PATH is required")
	}

	c.BaseURL = os.Getenv("BASE_URL")
	if c.BaseURL == "" {
		return nil, fmt.Errorf("BASE_URL is required")
	}

	c.SMTPHost = os.Getenv("SMTP_HOST")
	if c.SMTPHost == "" {
		return nil, fmt.Errorf("SMTP_HOST is required")
	}

	smtpPortStr := os.Getenv("SMTP_PORT")
	if smtpPortStr == "" {
		return nil, fmt.Errorf("SMTP_PORT is required")
	}
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil || smtpPort <= 0 {
		return nil, fmt.Errorf("invalid SMTP_PORT: %q", smtpPortStr)
	}
	c.SMTPPort = smtpPort

	c.SMTPUser = os.Getenv("SMTP_USER")
	c.SMTPPass = os.Getenv("SMTP_PASS")
	c.SMTPFrom = os.Getenv("SMTP_FROM")
	if c.SMTPFrom == "" {
		return nil, fmt.Errorf("SMTP_FROM is required")
	}

	c.SiteTargetURL = os.Getenv("SITE_TARGET_URL")
	if c.SiteTargetURL == "" {
		return nil, fmt.Errorf("SITE_TARGET_URL is required")
	}

	headlessStr := getenvDefault("PLAYWRIGHT_HEADLESS", "true")
	c.PlaywrightHeadless = headlessStr != "false" && headlessStr != "0"

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

package automation

import "time"

// SiteSelectors defines all CSS/XPath selectors needed to drive the
// third-party parking registration site. All site-specific knowledge should
// be captured here so that changes to the site only require updating this
// configuration.
type SiteSelectors struct {
	ApartmentName   string
	LicensePlate    string
	VehicleMake     string
	VehicleModel    string
	ResidentName    string
	UnitNumber      string
	VisitorName     string
	ResidentEmail   string
	SubmitButton    string
	SuccessSelector string // element that indicates successful registration
}

// WaitStrategy configures how the automation waits for the site to respond
// after submitting the form.
type WaitStrategy struct {
	// OverallTimeout is the maximum time allowed for a single registration
	// attempt, including navigation, form fill, and confirmation.
	OverallTimeout time.Duration

	// AfterSubmitWait is an additional wait time after clicking submit,
	// used when the site does not expose a clear success element.
	AfterSubmitWait time.Duration
}

// Config bundles all configuration required by the automation service.
// SITE_TARGET_URL and PLAYWRIGHT_HEADLESS are provided via environment and
// mapped into this structure by the caller (e.g. main or server wiring).
type Config struct {
	TargetURL string
	Selectors SiteSelectors
	Wait      WaitStrategy
	Headless  bool
}

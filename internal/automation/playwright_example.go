package automation

import (
	"context"
	"fmt"
	"time"

	"vistor-parking-automation-vrr/internal/models"

	"github.com/playwright-community/playwright-go"
)

type playwrightService struct {
	cfg     Config
	pw      *playwright.Playwright
	browser playwright.Browser
}

// NewPlaywrightService initializes Playwright and launches a browser instance
// that can be reused across registration calls.
func NewPlaywrightService(cfg Config) (Service, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("playwright.Run failed: %w", err)
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(cfg.Headless),
	})
	if err != nil {
		_ = pw.Stop()
		return nil, fmt.Errorf("browser launch failed: %w", err)
	}

	return &playwrightService{
		cfg:     cfg,
		pw:      pw,
		browser: browser,
	}, nil
}

func (s *playwrightService) RegisterVisitor(ctx context.Context, vehicle models.Vehicle) (*Result, error) {
	timeout := s.cfg.Wait.OverallTimeout
	if timeout <= 0 {
		timeout = 2 * time.Minute
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	logs := []string{}
	log := func(msg string) { logs = append(logs, msg) }

	log("Starting registration for vehicle: " + vehicle.VehicleMake + " " + vehicle.VehicleModel)

	browserCtx, err := s.browser.NewContext()
	if err != nil {
		log("browser context error: " + err.Error())
		return &Result{Success: false, Error: "browser context error", Logs: logs}, err
	}
	defer browserCtx.Close()

	page, err := browserCtx.NewPage()
	if err != nil {
		log("page create error: " + err.Error())
		return &Result{Success: false, Error: "page create error", Logs: logs}, err
	}
	defer page.Close()

	log("Navigating to: " + s.cfg.TargetURL)
	if _, err := page.Goto(s.cfg.TargetURL); err != nil {
		log("goto error: " + err.Error())
		return &Result{Success: false, Error: "navigation error", Logs: logs}, err
	}

	// Fill form fields using selectors from Config.
	if s.cfg.Selectors.LicensePlate != "" {
		log("Filling license plate: " + vehicle.LicensePlate)
		if err := page.Fill(s.cfg.Selectors.LicensePlate, vehicle.LicensePlate); err != nil {
			log("fill license plate error: " + err.Error())
			return &Result{Success: false, Error: "fill license plate", Logs: logs}, err
		}
	}

	if s.cfg.Selectors.VehicleMake != "" {
		log("Filling vehicle make: " + vehicle.VehicleMake)
		if err := page.Fill(s.cfg.Selectors.VehicleMake, vehicle.VehicleMake); err != nil {
			log("fill vehicle make error: " + err.Error())
		}
	}

	if s.cfg.Selectors.VehicleModel != "" {
		log("Filling vehicle model: " + vehicle.VehicleModel)
		if err := page.Fill(s.cfg.Selectors.VehicleModel, vehicle.VehicleModel); err != nil {
			log("fill vehicle model error: " + err.Error())
		}
	}

	if s.cfg.Selectors.ResidentName != "" {
		log("Filling resident name: " + vehicle.ResidentName)
		if err := page.Fill(s.cfg.Selectors.ResidentName, vehicle.ResidentName); err != nil {
			log("fill resident name error: " + err.Error())
		}
	}

	if s.cfg.Selectors.UnitNumber != "" {
		log("Filling unit number: " + vehicle.UnitNumber)
		if err := page.Fill(s.cfg.Selectors.UnitNumber, vehicle.UnitNumber); err != nil {
			log("fill unit number error: " + err.Error())
		}
	}

	if s.cfg.Selectors.VisitorName != "" {
		log("Filling visitor name: " + vehicle.VisitorName)
		if err := page.Fill(s.cfg.Selectors.VisitorName, vehicle.VisitorName); err != nil {
			log("fill visitor name error: " + err.Error())
		}
	}

	if s.cfg.Selectors.ResidentEmail != "" {
		log("Filling confirmation email: " + vehicle.ConfirmationEmail)
		if err := page.Fill(s.cfg.Selectors.ResidentEmail, vehicle.ConfirmationEmail); err != nil {
			log("fill confirmation email error: " + err.Error())
		}
	}

	if s.cfg.Selectors.ApartmentName != "" {
		log("Filling apartment name: " + vehicle.ApartmentName)
		if err := page.Fill(s.cfg.Selectors.ApartmentName, vehicle.ApartmentName); err != nil {
			log("fill apartment name error: " + err.Error())
		}
	}

	if s.cfg.Selectors.SubmitButton != "" {
		log("Clicking submit button")
		if err := page.Click(s.cfg.Selectors.SubmitButton); err != nil {
			log("click submit error: " + err.Error())
			return &Result{Success: false, Error: "submit error", Logs: logs}, err
		}
	}

	// Wait for success indication or a fixed delay.
	if s.cfg.Selectors.SuccessSelector != "" {
		log("Waiting for success selector: " + s.cfg.Selectors.SuccessSelector)
		if _, err := page.WaitForSelector(s.cfg.Selectors.SuccessSelector); err != nil {
			log("wait for success selector error: " + err.Error())
		}
	} else if s.cfg.Wait.AfterSubmitWait > 0 {
		log(fmt.Sprintf("Waiting %v after submit", s.cfg.Wait.AfterSubmitWait))
		select {
		case <-ctx.Done():
			return &Result{Success: false, Error: "timeout waiting after submit", Logs: logs}, ctx.Err()
		case <-time.After(s.cfg.Wait.AfterSubmitWait):
		}
	}

	log("Registration completed successfully")
	return &Result{Success: true, Logs: logs}, nil
}

func (s *playwrightService) Close() error {
	if s.browser != nil {
		if err := s.browser.Close(); err != nil {
			return err
		}
	}
	if s.pw != nil {
		if err := s.pw.Stop(); err != nil {
			return err
		}
	}
	return nil
}

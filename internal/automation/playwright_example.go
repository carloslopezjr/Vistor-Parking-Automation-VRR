package automation

import (
	"context"
	"fmt"
	"log"
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
	// Install browsers and system dependencies if needed
	if err := playwright.Install(); err != nil {
		return nil, fmt.Errorf("playwright.Install failed: %w", err)
	}

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
	log := func(msg string) {
		logs = append(logs, msg)
		log.Println("[AUTOMATION] " + msg)
	}

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

	// ==== PAGE 1: Apartment Selection ====
	// Step 1: Fill apartment name
	if s.cfg.Selectors.ApartmentName != "" {
		log("Waiting for apartment name selector: " + s.cfg.Selectors.ApartmentName)
		if _, err := page.WaitForSelector(s.cfg.Selectors.ApartmentName, playwright.PageWaitForSelectorOptions{
			Timeout: playwright.Float(100000),
		}); err != nil {
			log("apartment name selector not found (5s timeout): " + err.Error())
			title, _ := page.Title()
			log("Page title: " + title)
			log("Page URL: " + page.URL())
			return &Result{Success: false, Error: "apartment name field not found", Logs: logs}, err
		}

		log("Found apartment name selector, filling: " + vehicle.ApartmentName)
		if err := page.Fill(s.cfg.Selectors.ApartmentName, vehicle.ApartmentName); err != nil {
			log("fill apartment name error: " + err.Error())
			return &Result{Success: false, Error: "fill apartment name", Logs: logs}, err
		}
		select {
		case <-time.After(500 * time.Millisecond):
		}
		log("Successfully filled apartment name")

		// Step 2: Wait 1 second for button to appear
		log("Waiting 1 second for apartment confirm button to load...")
		select {
		case <-time.After(1 * time.Second):
		}

		// Step 3: Click apartment confirm button
		if s.cfg.Selectors.ApartmentConfirmButton != "" {
			log("Waiting for apartment confirm button: " + s.cfg.Selectors.ApartmentConfirmButton)
			if _, err := page.WaitForSelector(s.cfg.Selectors.ApartmentConfirmButton, playwright.PageWaitForSelectorOptions{
				Timeout: playwright.Float(5000),
			}); err != nil {
				log("apartment confirm button not found (5s timeout): " + err.Error())
				return &Result{Success: false, Error: "apartment confirm button not found", Logs: logs}, err
			}

			log("Clicking apartment confirm button")
			if err := page.Click(s.cfg.Selectors.ApartmentConfirmButton); err != nil {
				log("click apartment confirm button error: " + err.Error())
				return &Result{Success: false, Error: "click apartment confirm button", Logs: logs}, err
			}
			log("Successfully clicked apartment confirm button")

			// Wait for page to navigate to next page (license plate page)
			log("Waiting for navigation to license plate page...")
			select {
			case <-time.After(2 * time.Second):
			}
		}
	}

	// ==== PAGE 2: License Plate & Vehicle Details ====
	// Fill license plate
	if s.cfg.Selectors.LicensePlate != "" {
		log("Waiting for license plate selector: " + s.cfg.Selectors.LicensePlate)
		if _, err := page.WaitForSelector(s.cfg.Selectors.LicensePlate, playwright.PageWaitForSelectorOptions{
			Timeout: playwright.Float(5000),
		}); err != nil {
			log("license plate selector not found (5s timeout): " + err.Error())
			title, _ := page.Title()
			log("Page title: " + title)
			log("Page URL: " + page.URL())
			return &Result{Success: false, Error: "license plate field not found", Logs: logs}, err
		}

		log("Found license plate selector, filling: " + vehicle.LicensePlate)
		if err := page.Fill(s.cfg.Selectors.LicensePlate, vehicle.LicensePlate); err != nil {
			log("fill license plate error: " + err.Error())
			return &Result{Success: false, Error: "fill license plate", Logs: logs}, err
		}
		select {
		case <-time.After(500 * time.Millisecond):
		}
		log("Successfully filled license plate")
	}

	// Fill confirmed license plate
	if s.cfg.Selectors.ConfirmedLicensePlate != "" {
		log("Waiting for confirmed license plate selector: " + s.cfg.Selectors.ConfirmedLicensePlate)
		if _, err := page.WaitForSelector(s.cfg.Selectors.ConfirmedLicensePlate, playwright.PageWaitForSelectorOptions{
			Timeout: playwright.Float(5000),
		}); err != nil {
			log("confirmed license plate selector not found (5s timeout): " + err.Error())
		} else {
			log("Found confirmed license plate selector, filling: " + vehicle.ConfirmedLicensePlate)
			if err := page.Fill(s.cfg.Selectors.ConfirmedLicensePlate, vehicle.ConfirmedLicensePlate); err != nil {
				log("fill confirmed license plate error: " + err.Error())
			} else {
				select {
				case <-time.After(500 * time.Millisecond):
				}
				log("Successfully filled confirmed license plate")
			}
		}
	}

	// ==== Continue with remaining fields ====
	if s.cfg.Selectors.VehicleMake != "" {
		log("Filling vehicle make: " + vehicle.VehicleMake)
		if err := page.Fill(s.cfg.Selectors.VehicleMake, vehicle.VehicleMake); err != nil {
			log("fill vehicle make error: " + err.Error())
		} else {
			select {
			case <-time.After(500 * time.Millisecond):
			}
		}
	}

	if s.cfg.Selectors.VehicleModel != "" {
		log("Filling vehicle model: " + vehicle.VehicleModel)
		if err := page.Fill(s.cfg.Selectors.VehicleModel, vehicle.VehicleModel); err != nil {
			log("fill vehicle model error: " + err.Error())
		} else {
			select {
			case <-time.After(500 * time.Millisecond):
			}
		}
	}

	if s.cfg.Selectors.ResidentName != "" {
		log("Filling resident name: " + vehicle.ResidentName)
		if err := page.Fill(s.cfg.Selectors.ResidentName, vehicle.ResidentName); err != nil {
			log("fill resident name error: " + err.Error())
		} else {
			select {
			case <-time.After(500 * time.Millisecond):
			}
		}
	}

	if s.cfg.Selectors.UnitNumberButton != "" {
		log("Clicking unit number button: " + s.cfg.Selectors.UnitNumberButton)
		if err := page.Click(s.cfg.Selectors.UnitNumberButton); err != nil {
			log("click unit number button error: " + err.Error())
		} else {
			log("Successfully clicked unit number button")

			// Step 2: Wait for unit number input field to appear in dialog
			if s.cfg.Selectors.UnitNumberInput != "" {
				log("Waiting for unit number input field: " + s.cfg.Selectors.UnitNumberInput)
				if _, err := page.WaitForSelector(s.cfg.Selectors.UnitNumberInput, playwright.PageWaitForSelectorOptions{
					Timeout: playwright.Float(5000),
				}); err != nil {
					log("unit number input not found (5s timeout): " + err.Error())
				} else {
					log("Found unit number input field, filling: " + vehicle.UnitNumber)
					if err := page.Fill(s.cfg.Selectors.UnitNumberInput, vehicle.UnitNumber); err != nil {
						log("fill unit number error: " + err.Error())
					} else {
						select {
						case <-time.After(500 * time.Millisecond):
						}
						log("Successfully filled unit number")

						// Step 4: Click unit number confirm button
						if s.cfg.Selectors.UnitNumberConfirmButton != "" {
							log("Clicking unit number confirm button: " + s.cfg.Selectors.UnitNumberConfirmButton)
							if err := page.Click(s.cfg.Selectors.UnitNumberConfirmButton); err != nil {
								log("click unit number confirm button error: " + err.Error())
							} else {
								log("Successfully clicked unit number confirm button")
								// Wait for dialog to close
								select {
								case <-time.After(1 * time.Second):
								}
							}
						}
					}
				}
			}
		}
	}

	if s.cfg.Selectors.VisitorName != "" {
		log("Filling visitor name: " + vehicle.VisitorName)
		if err := page.Fill(s.cfg.Selectors.VisitorName, vehicle.VisitorName); err != nil {
			log("fill visitor name error: " + err.Error())
		} else {
			select {
			case <-time.After(500 * time.Millisecond):
			}
		}
	}

	if s.cfg.Selectors.ResidentEmail != "" {
		log("Filling confirmation email: " + vehicle.ConfirmationEmail)
		if err := page.Fill(s.cfg.Selectors.ResidentEmail, vehicle.ConfirmationEmail); err != nil {
			log("fill confirmation email error: " + err.Error())
		} else {
			select {
			case <-time.After(500 * time.Millisecond):
			}
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

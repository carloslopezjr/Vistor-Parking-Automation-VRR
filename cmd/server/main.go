package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vistor-parking-automation-vrr/internal/automation"
	"vistor-parking-automation-vrr/internal/config"
	csvloader "vistor-parking-automation-vrr/internal/csv"
	apphttp "vistor-parking-automation-vrr/internal/http"
)

func main() {
	logger := log.Default()

	// Best-effort load of .env (or .env.local) before reading configuration.
	_ = config.LoadDotEnv(".env")
	_ = config.LoadDotEnv(".env.local")

	// Load and validate configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("config error: %v", err)
	}

	// Load vehicles from CSV
	logger.Printf("Loading vehicles from %s", cfg.VehiclesCSVPath)
	vehicles, err := csvloader.LoadVehicles(cfg.VehiclesCSVPath)
	if err != nil {
		logger.Fatalf("csv load error: %v", err)
	}
	logger.Printf("Loaded %d vehicles", len(vehicles))

	// Initialize Playwright automation service
	logger.Printf("Initializing Playwright automation")
	automator, err := automation.NewPlaywrightService(automation.Config{
		TargetURL: cfg.PlaywrightTargetURL,
		Selectors: cfg.PlaywrightSelectors,
		Wait: automation.WaitStrategy{
			OverallTimeout:  cfg.PlaywrightTimeout,
			AfterSubmitWait: cfg.PlaywrightWaitAfterSubmit,
		},
		Headless: cfg.PlaywrightHeadless,
	})
	if err != nil {
		logger.Fatalf("playwright init error: %v", err)
	}
	defer func() {
		if err := automator.Close(); err != nil {
			logger.Printf("error closing playwright: %v", err)
		}
	}()

	// Create HTTP server
	server, err := apphttp.NewServer(vehicles, cfg, automator)
	if err != nil {
		logger.Fatalf("http server init error: %v", err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: server.Handler(),
	}

	go func() {
		logger.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server error: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Printf("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Printf("server shutdown error: %v", err)
	}
}

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
	"vistor-parking-automation-vrr/internal/db"
	apphttp "vistor-parking-automation-vrr/internal/http"
	"vistor-parking-automation-vrr/internal/jobs"
	"vistor-parking-automation-vrr/internal/mailer"
	"vistor-parking-automation-vrr/internal/scheduler"
	"vistor-parking-automation-vrr/internal/tokens"
)

func main() {
	logger := log.Default()

	// Best-effort load of .env (or .env.local) before reading configuration.
	_ = config.LoadDotEnv(".env")
	_ = config.LoadDotEnv(".env.local")

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("config error: %v", err)
	}

	database, err := db.Open(cfg.DatabasePath)
	if err != nil {
		logger.Fatalf("db open error: %v", err)
	}
	defer database.Close()

	ctx := context.Background()
	if err := db.Migrate(ctx, database); err != nil {
		logger.Fatalf("db migrate error: %v", err)
	}

	// Core services
	tokenSvc := tokens.NewService(database, tokens.Config{TTL: 48 * time.Hour})
	mailSvc := mailer.New(mailer.Config{
		Host: cfg.SMTPHost,
		Port: cfg.SMTPPort,
		User: cfg.SMTPUser,
		Pass: cfg.SMTPPass,
		From: cfg.SMTPFrom,
	})

	// Placeholder automator; real implementation wired later.
	automator := automation.NewNoopService()

	// Jobs and scheduler
	jobsSvc := jobs.NewService(database, nil) // reminder handler wired later
	sched := scheduler.New(database, jobsSvc, logger, scheduler.Config{})

	server, err := apphttp.NewServer(database, cfg, jobsSvc, mailSvc, tokenSvc, automator)
	if err != nil {
		logger.Fatalf("http server init error: %v", err)
	}

	// Run scheduler in background.
	go func() {
		if err := sched.Run(context.Background()); err != nil && err != context.Canceled {
			logger.Printf("scheduler exited: %v", err)
		}
	}()

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

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Printf("server shutdown error: %v", err)
	}
}

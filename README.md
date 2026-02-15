# Vistor-Parking-Automation-VRR

Automated sign up for frequent visitors using vrrparking.com

## Visitor Parking Automation (MVP)

This Go web application automates re-registering visitor parking by auto-filling a third-party website form. It is intended for authorized use only.

### Features

- SQLite-backed storage for parking profiles
- Mobile-first web UI with large profile buttons
- One-click reminder emails with secure, single-use tokens
- DB-backed scheduler for reminder jobs (12h/20h/23h)
- Pluggable automation service (Playwright-based example included)

### Configuration

All configuration is via environment variables:

- `PORT`
- `DATABASE_PATH`
- `BASE_URL`
- `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS`, `SMTP_FROM`
- `SITE_TARGET_URL`
- `PLAYWRIGHT_HEADLESS`
- `BASIC_AUTH_USER`, `BASIC_AUTH_PASS` (optional)

### Running Locally

```bash
go run ./cmd/server
```

Ensure `DATABASE_PATH` points to a writable location.

### Docker

Build and run:

```bash
docker build -t visitor-parking .
docker run --rm -p 8080:8080 -e DATABASE_PATH=/data/app.db visitor-parking
```

### Notes

- Automation is wired via an interface; a commented Playwright-based implementation is provided in `internal/automation/playwright_example.go` and may require adjusting to match the chosen Playwright Go SDK version.

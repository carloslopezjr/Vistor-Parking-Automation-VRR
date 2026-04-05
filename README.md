# Visitor Parking Automation - VRR

Automated visitor parking registration for vrrparking.com using CSV and Playwright browser automation.

## Overview

This lightweight Go web application automates visitor parking registration by loading vehicles from a CSV file and using Playwright to auto-fill and submit the registration form. It's designed for authorized personal use to streamline frequent visitor registration.

**Key Principles:**

- Simple CSV-based configuration (no database)
- One-click automation with real-time logging
- Headless browser automation with Playwright
- Request throttling (500ms delays) to prevent rate limiting
- No background jobs or complex scheduling

## Features

- **CSV Vehicle Management** - Load vehicle details from `vehicles.csv`
- **Web UI** - Simple vehicle selector with one-click registration button
- **Playwright Automation** - Headless browser automation with automatic form filling
- **Real-time Logging** - Terminal and UI logging of all automation steps
- **Request Throttling** - 500ms delays between form field fills to prevent bans
- **XPath & CSS Selectors** - Flexible form element targeting via environment configuration

## Quick Start

### 1. Setup

Copy the example CSV:

```bash
cp vehicles.csv.example vehicles.csv
```

Edit `vehicles.csv` with your vehicle information:

```csv
apartment_name,license_plate,confirm_license_plate,vehicle_make,vehicle_model,resident_name,unit_number,visitor_name,confirmation_email
The Merle,ABC1234,ABC1234,Toyota,Camry,John Smith,101,Jane Doe,john@example.com
```

### 2. Configure

Copy `.env.example` to `.env` and update the XPath selectors if the form structure changes:

```bash
cp .env.example .env
```

Edit `.env` and update `PLAYWRIGHT_SELECTOR_*` values if needed.

### 3. Run

```bash
go run ./cmd/server
```

Visit `http://localhost:8080` and select a vehicle to auto-register.

## Configuration

All configuration via environment variables in `.env`:

**Server:**

- `PORT` - HTTP server port (default: 8080)
- `BASE_URL` - Application base URL
- `BASIC_AUTH_USER` / `BASIC_AUTH_PASS` - Optional basic auth

**CSV:**

- `VEHICLES_CSV_PATH` - Path to CSV file (default: ./vehicles.csv)

**Playwright:**

- `PLAYWRIGHT_TARGET_URL` - vrrparking.com form URL
- `PLAYWRIGHT_HEADLESS` - false to see browser (default: true)
- `PLAYWRIGHT_TIMEOUT` - Max registration time (default: 120s)
- `PLAYWRIGHT_WAIT_AFTER_SUBMIT` - Wait after submit (default: 5s)
- `PLAYWRIGHT_SELECTOR_*` - Form field XPath/CSS selectors

## Security Notes

⚠️ **Local Use Only:**

- `vehicles.csv` contains personal information (license plates, addresses, emails)
- It's in `.gitignore` to prevent accidental commits
- Never commit real vehicle data to version control
- For production use, implement proper credential management

## Architecture

- **`cmd/server/main.go`** - Application entry point
- **`internal/automation/`** - Playwright-based browser automation
- **`internal/csv/`** - CSV vehicle loader
- **`internal/http/`** - Web server and handlers
- **`internal/config/`** - Environment configuration
- **`web/templates/`** - HTML/JavaScript UI

## Docker

Build and run:

```bash
docker build -t visitor-parking .
docker run --rm -p 8080:8080 \
  -e PLAYWRIGHT_HEADLESS=false \
  -v $(pwd)/vehicles.csv:/app/vehicles.csv \
  -v $(pwd)/.env:/app/.env \
  visitor-parking
```

**Note:** Docker image will install Playwright browsers on first run (adds ~500MB).

## Development

Run tests:

```bash
go test ./...
```

Build binary:

```bash
go build -o visitor-parking ./cmd/server
```

## License

For authorized personal use only. Ensure compliance with vrrparking.com terms of service.

## Troubleshooting

**Form fields not filling?**

- Inspect the website with browser dev tools
- Update XPath selectors in `.env`
- Check logs for "selector not found" errors

**Getting rate-limited/banned?**

- Increase `PLAYWRIGHT_WAIT_AFTER_SUBMIT`
- Add manual delays by editing selector wait times
- Space out registrations with wait time

# Multi-stage build for Visitor Parking Automation
# The application uses Playwright for browser automation, which requires system dependencies.
FROM golang:1.22-alpine AS build

WORKDIR /app

# Install build dependencies (required for go-sqlite3 if needed, and other CGO dependencies)
RUN apk add --no-cache build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# Runtime image with Playwright dependencies
FROM alpine:3.19

# Install Playwright/browser runtime dependencies
# These are needed for Chromium browser to run in headless mode
RUN apk add --no-cache \
    ca-certificates \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ttf-dejavu

# Create non-root user for security
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy binary and web assets from build stage
COPY --from=build /app/server /app/server
COPY --from=build /app/web /app/web

# Copy example files for reference
COPY vehicles.csv.example /app/vehicles.csv.example
COPY .env.example /app/.env.example

# Set default environment variables
ENV PORT=8080
ENV VEHICLES_CSV_PATH=/app/vehicles.csv
ENV PLAYWRIGHT_HEADLESS=true

EXPOSE 8080

USER appuser

CMD ["/app/server"]

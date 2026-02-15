# Minimal production-minded Dockerfile for the Go web app
FROM golang:1.22-alpine AS build

WORKDIR /app

RUN apk add --no-cache build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM alpine:3.19

RUN adduser -D -g '' appuser

WORKDIR /app
COPY --from=build /app/server /app/server
COPY --from=build /app/web /app/web

# Create a writable data directory for SQLite and mark it as a volume.
RUN mkdir -p /data && chmod 777 /data
VOLUME ["/data"]

USER appuser

ENV PORT=8080
EXPOSE 8080

CMD ["/app/server"]

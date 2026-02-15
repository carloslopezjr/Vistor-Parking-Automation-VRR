package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logging wraps an http.Handler and logs basic request information.
func Logging(next http.Handler) http.Handler {
	logger := log.Default()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		logger.Printf("http request method=%s path=%s status=%d duration=%s", r.Method, r.URL.Path, rw.status, time.Since(start))
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

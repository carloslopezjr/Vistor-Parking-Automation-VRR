package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// BasicAuthConfig holds optional basic authentication credentials.
type BasicAuthConfig struct {
	User string
	Pass string
}

// BasicAuth wraps an http.Handler with HTTP Basic authentication when
// credentials are configured. If User is empty, the handler is returned
// unchanged.
func BasicAuth(next http.Handler, cfg BasicAuthConfig) http.Handler {
	if cfg.User == "" {
		return next
	}

	realm := "Restricted"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if h == "" || !strings.HasPrefix(h, "Basic ") {
			w.Header().Set("WWW-Authenticate", "Basic realm="+realm)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(h, "Basic "))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(string(payload), ":", 2)
		if len(parts) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, pass := parts[0], parts[1]
		if subtleConstantTimeCompare(user, cfg.User) && subtleConstantTimeCompare(pass, cfg.Pass) {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("WWW-Authenticate", "Basic realm="+realm)
		w.WriteHeader(http.StatusUnauthorized)
	})
}

// subtleConstantTimeCompare performs a constant-time comparison to avoid
// leaking timing information about credentials.
func subtleConstantTimeCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	var v byte
	for i := 0; i < len(a); i++ {
		v |= a[i] ^ b[i]
	}
	return v == 0
}

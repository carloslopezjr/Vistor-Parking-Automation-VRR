package http

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"vistor-parking-automation-vrr/internal/automation"
	appcfg "vistor-parking-automation-vrr/internal/config"
	"vistor-parking-automation-vrr/internal/http/handlers"
	"vistor-parking-automation-vrr/internal/http/middleware"
	"vistor-parking-automation-vrr/internal/jobs"
	"vistor-parking-automation-vrr/internal/mailer"
	"vistor-parking-automation-vrr/internal/store"
	"vistor-parking-automation-vrr/internal/tokens"
)

// Server wraps the HTTP router and dependencies.
type Server struct {
	handler http.Handler
}

// NewServer constructs the HTTP server with all routes configured.
func NewServer(db *sql.DB, cfg *appcfg.Config, jobsSvc jobs.Service, mailSvc mailer.Service, tokenSvc tokens.Service, automator automation.Service) (*Server, error) {
	mux := http.NewServeMux()

	profiles := store.NewProfileStore(db)
	logsStore := store.NewLogStore(db)

	tpl, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		return nil, err
	}

	homeHandler := &handlers.HomeHandler{
		Profiles: profiles,
		Logs:     logsStore,
		Tpl:      tpl,
	}

	profilesHandler := &handlers.ProfilesHandler{
		Profiles: profiles,
		Tpl:      tpl,
	}

	tokenHandler := &handlers.TokenHandler{
		Tokens:    tokenSvc,
		Profiles:  profiles,
		Jobs:      jobsSvc,
		Mailer:    mailSvc,
		BaseURL:   cfg.BaseURL,
		Automator: automator,
	}

	mux.Handle("/", homeHandler)
	mux.Handle("/r/", tokenHandler)
	mux.Handle("/profiles", profilesHandler)
	mux.Handle("/profiles/", profilesHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	root := http.Handler(mux)
	root = middleware.Recovery(root)
	root = middleware.Logging(root)
	root = middleware.BasicAuth(root, middleware.BasicAuthConfig{
		User: cfg.BasicAuthUser,
		Pass: cfg.BasicAuthPass,
	})

	_ = log.Default() // reserved for future structured logging integration

	return &Server{handler: root}, nil
}

// Handler returns the underlying http.Handler.
func (s *Server) Handler() http.Handler {
	return s.handler
}

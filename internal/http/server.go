package http

import (
	"html/template"
	"net/http"

	"vistor-parking-automation-vrr/internal/automation"
	appcfg "vistor-parking-automation-vrr/internal/config"
	"vistor-parking-automation-vrr/internal/http/handlers"
	"vistor-parking-automation-vrr/internal/http/middleware"
	"vistor-parking-automation-vrr/internal/models"
)

// Server wraps the HTTP router and dependencies.
type Server struct {
	handler http.Handler
}

// NewServer constructs the HTTP server with all routes configured.
func NewServer(vehicles []models.Vehicle, cfg *appcfg.Config, automator automation.Service) (*Server, error) {
	mux := http.NewServeMux()

	tpl, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		return nil, err
	}

	homeHandler := &handlers.HomeHandler{
		Vehicles: vehicles,
		Tpl:      tpl,
	}

	registerHandler := &handlers.RegisterHandler{
		Vehicles:  vehicles,
		Automator: automator,
	}

	mux.Handle("/", homeHandler)
	mux.Handle("/register", registerHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	root := http.Handler(mux)
	root = middleware.Recovery(root)
	root = middleware.Logging(root)
	root = middleware.BasicAuth(root, middleware.BasicAuthConfig{
		User: cfg.BasicAuthUser,
		Pass: cfg.BasicAuthPass,
	})

	return &Server{handler: root}, nil
}

// Handler returns the underlying http.Handler.
func (s *Server) Handler() http.Handler {
	return s.handler
}

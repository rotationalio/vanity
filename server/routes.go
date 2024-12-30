package server

import (
	"io/fs"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rotationalio/vanity/logger"
	"github.com/rotationalio/vanity/server/middleware"
)

// Sets up the server's middleware and routes.
func (s *Server) setupRoutes() (err error) {
	middleware := []middleware.Middleware{
		logger.HTTPLogger("vanity"),
		s.Maintenance(),
	}

	// Kubernetes liveness probes added before middleware.
	s.router.GET("/healthz", s.Healthz)
	s.router.GET("/livez", s.Healthz)
	s.router.GET("/readyz", s.Readyz)

	// API Routes
	// Status/Heartbeat endpoint
	s.addRoute(http.MethodGet, "/v1/status", s.Status, middleware...)

	// Application Routes
	static, _ := fs.Sub(content, "static")
	s.router.ServeFiles("/static/*filepath", http.FS(static))
	s.addRoute(http.MethodGet, "/", s.HomePage(), middleware...)

	// Golang Vanity Handling

	return nil
}

func (s *Server) addRoute(method, path string, h httprouter.Handle, m ...middleware.Middleware) {
	s.router.Handle(method, path, middleware.Chain(h, m...))
}

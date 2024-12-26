package server

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rotationalio/vanity"
	"github.com/rotationalio/vanity/server/middleware"
	"github.com/rotationalio/vanity/server/render"
)

// If the server is in maintenance mode, aborts the current request and renders the
// maintenance mode page instead. Returns nil if not in maintenance mode.
func (s *Server) Maintenance() middleware.Middleware {
	if s.conf.Maintenance {
		return func(next httprouter.Handle) httprouter.Handle {
			return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				render.JSON(http.StatusServiceUnavailable, w, &StatusReply{
					Status:  "maintenance",
					Version: vanity.Version(),
					Uptime:  time.Since(s.started).String(),
				})
			}
		}
	}
	return nil
}

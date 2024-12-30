package server

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//go:embed all:templates
//go:embed all:static
var content embed.FS

func (s *Server) HomePage() httprouter.Handle {
	// Compile the template for serving the home page.
	templates, _ := fs.Sub(content, "templates")
	index := template.Must(template.ParseFS(templates, "*.html"))
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if err := index.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

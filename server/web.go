package server

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rotationalio/vanity"
)

//go:embed all:templates
//go:embed all:static
var content embed.FS

type homeTableItem struct {
	Module string
	Import string
	Source string
	GoDoc  string
}

func (s *Server) HomePage(pkgs []*vanity.GoPackage) httprouter.Handle {
	// Compile the template for serving the home page.
	templates, _ := fs.Sub(content, "templates")
	index := template.Must(template.ParseFS(templates, "*.html"))

	// Compile the package information for serving the modules table
	table := make([]homeTableItem, 0, len(pkgs))
	for _, pkg := range pkgs {
		table = append(table, homeTableItem{
			Module: pkg.Module,
			Import: pkg.Import(),
			Source: pkg.Repository,
			GoDoc:  pkg.Redirect(),
		})
	}

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if err := index.ExecuteTemplate(w, "index.html", table); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

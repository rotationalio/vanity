package server

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rotationalio/vanity"
)

func Vanity(pkg *vanity.GoPackage) httprouter.Handle {
	// Compile the template for serving the vanity url.
	templates, _ := fs.Sub(content, "templates")
	index := template.Must(template.ParseFS(templates, "*.html"))

	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Construct a request specific response
		data := pkg.WithRequest(r)

		// Issue an HTTP redirect if this is definitely a browser.
		if r.FormValue("go-get") != "1" {
			http.Redirect(w, r, data.Redirect(), http.StatusTemporaryRedirect)
			return
		}

		// Write go-import and go-source meta tags to response.
		w.Header().Set("Cache-Control", "public, max-age=300")
		if err := index.ExecuteTemplate(w, "vanity.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

package render

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Header keys for http requests and responses
const (
	Accept      = "Accept"
	ContentType = "Content-Type"
)

// Content type values
var (
	plainContentType = "text/plain"
	jsonContentType  = "application/json"
)

func Text(code int, w http.ResponseWriter, text string) error {
	w.Header().Set(ContentType, plainContentType)
	w.WriteHeader(code)

	fmt.Fprintln(w, text)
	return nil
}

func Textf(code int, w http.ResponseWriter, text string, a ...any) error {
	w.Header().Set(ContentType, plainContentType)
	w.WriteHeader(code)

	fmt.Fprintf(w, text, a...)
	return nil
}

// JSON marshals the given interface object and writes it with the correct ContentType.
func JSON(code int, w http.ResponseWriter, obj any) error {
	w.Header().Set(ContentType, jsonContentType)
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(obj); err != nil {
		return err
	}
	return nil
}

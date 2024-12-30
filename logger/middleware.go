package logger

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rotationalio/vanity"
	"github.com/rotationalio/vanity/server/middleware"
	"github.com/rs/zerolog/log"
)

func HTTPLogger(server string) middleware.Middleware {
	version := vanity.Version()
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			// Before the request
			started := time.Now()
			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path = path + "?" + r.URL.RawQuery
			}

			// Wrap the response writer for the logger
			rw := &LogWriter{ResponseWriter: w}

			// Handle the request
			next(rw, r, p)

			// After the request
			status := rw.Status()
			logctx := log.With().
				Str("path", path).
				Str("ser_name", server).
				Str("version", version).
				Str("method", r.Method).
				Dur("resp_time", time.Since(started)).
				Int("resp_bytes", rw.Size()).
				Int("status", status).
				Logger()

			switch {
			case status >= 400 && status < 500:
				logctx.Warn().Msgf("%s %s %s %d", server, r.Method, path, status)
			case status >= 500:
				logctx.Error().Msgf("%s %s %s %d", server, r.Method, path, status)
			default:
				logctx.Info().Msgf("%s %s %s %d", server, r.Method, path, status)
			}

		}
	}
}

type LogWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *LogWriter) WriteHeader(statusCode int) {
	if w.Written() {
		return
	}

	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *LogWriter) Write(b []byte) (int, error) {
	if !w.Written() {
		w.WriteHeader(http.StatusOK)
	}

	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

func (w *LogWriter) Written() bool {
	return w.status >= http.StatusOK || w.status == http.StatusSwitchingProtocols
}

func (w *LogWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *LogWriter) Status() int {
	return w.status
}

func (w *LogWriter) Size() int {
	return w.size
}

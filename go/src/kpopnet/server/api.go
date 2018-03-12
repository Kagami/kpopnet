// Server API handlers.
package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"kpopnet/db"
)

var (
	// Client errors.
	// TODO(Kagami): Use custom class to store error context.
	errInternal = fmt.Errorf("internal-error")
)

func serveData(w http.ResponseWriter, r *http.Request, data interface{}) {
	buf, err := json.Marshal(data)
	if err != nil {
		handle500(w, r, err)
		return
	}
	etag := fmt.Sprintf("\"%s\"", hashBytes(buf))
	if checkEtag(w, r, etag) {
		return
	}
	head := w.Header()
	head.Set("Cache-Control", "no-cache")
	head.Set("Content-Type", "application/json")
	head.Set("ETag", etag)
	w.Write(buf)
}

func serveError(w http.ResponseWriter, r *http.Request, err error, code int) {
	head := w.Header()
	head.Set("Cache-Control", "no-cache")
	head.Set("Content-Type", "application/json")
	w.WriteHeader(code)
	io.WriteString(w, fmt.Sprintf("{\"error\": \"%v\"}", err))
}

func handle500(w http.ResponseWriter, r *http.Request, err error) {
	logError(err)
	serveError(w, r, errInternal, 500)
}

func serveProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := db.GetProfiles()
	if err != nil {
		handle500(w, r, err)
		return
	}
	serveData(w, r, profiles)
}

package kpopnet

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	// Client errors.
	// TODO(Kagami): Use custom class to store error context.
	errInternal = errors.New("internal-error")
)

func setApiHeaders(w http.ResponseWriter) {
	head := w.Header()
	head.Set("Cache-Control", "no-cache")
	head.Set("Content-Type", "application/json")
}

func serveData(w http.ResponseWriter, r *http.Request, data []byte) {
	// Will need to hash on each request but this seems to be <1ms.
	etag := fmt.Sprintf("\"%s\"", hashBytes(data))
	if checkEtag(w, r, etag) {
		return
	}
	setApiHeaders(w)
	w.Header().Set("ETag", etag)
	w.Write(data)
}

func serveError(w http.ResponseWriter, r *http.Request, err error, code int) {
	setApiHeaders(w)
	w.WriteHeader(code)
	io.WriteString(w, fmt.Sprintf("{\"error\": \"%v\"}", err))
}

func handle500(w http.ResponseWriter, r *http.Request, err error) {
	logError(err)
	serveError(w, r, errInternal, 500)
}

func ServeProfiles(w http.ResponseWriter, r *http.Request) {
	v, err := cached(profileCacheKey, func() (v interface{}, err error) {
		ps, err := GetProfiles()
		if err != nil {
			return
		}
		// Takes ~5ms so better to store encoded.
		return json.Marshal(ps)
	})
	if err != nil {
		handle500(w, r, err)
		return
	}
	serveData(w, r, v.([]byte))
}

// Idol API is served by cutechan-compatible backend.
func serveIdolApi(w http.ResponseWriter, r *http.Request) {
	url := idolApi + "/api/idols/" + getParam(r, "path")
	http.Redirect(w, r, url, 302)
}

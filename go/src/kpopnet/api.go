package kpopnet

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	maxOverheadSize = int64(10 * 1024)
	maxFileSize     = int64(5 * 1024 * 1024)
	maxBodySize     = maxFileSize + maxOverheadSize

	// TODO(Kagami): Better error granularity?
	errInternal     = errors.New("Internal error")
	errParseForm    = errors.New("Error parsing form")
	errParseFile    = errors.New("Error parsing form file")
	errBadImage     = errors.New("Invalid image")
	errNoSingleFace = errors.New("Not a single face")
	errNoIdol       = errors.New("Cannot find idol")
)

func setMainHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache")
}

func setBodyHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func setApiHeaders(w http.ResponseWriter) {
	setMainHeaders(w)
	setBodyHeaders(w)
}

func serveData(w http.ResponseWriter, r *http.Request, data []byte) {
	setMainHeaders(w)
	etag := fmt.Sprintf("W/\"%s\"", hashBytes(data))
	w.Header().Set("ETag", etag)
	if checkEtag(w, r, etag) {
		return
	}
	setBodyHeaders(w)
	w.Write(data)
}

func serveJson(w http.ResponseWriter, r *http.Request, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		handle500(w, r, err)
		return
	}
	serveData(w, r, data)
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
	// TODO(Kagami): For some reason cached request is not fast enough.
	// TODO(Kagami): Use some trigger to invalidate cache.
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

func ServeRecognize(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	if err := r.ParseMultipartForm(0); err != nil {
		serveError(w, r, errParseForm, 400)
		return
	}
	fhs := r.MultipartForm.File["files[]"]
	if len(fhs) != 1 {
		serveError(w, r, errParseFile, 400)
		return
	}
	idolId, err := RequestRecognizeMultipart(fhs[0])
	switch err {
	case errParseFile:
		serveError(w, r, err, 400)
		return
	case errBadImage:
		serveError(w, r, err, 400)
		return
	case errNoIdol:
		serveError(w, r, err, 404)
		return
	case nil:
		// Do nothing.
	default:
		handle500(w, r, err)
		return
	}
	if idolId == nil {
		serveError(w, r, errNoSingleFace, 400)
		return
	}
	result := map[string]string{"id": *idolId}
	serveJson(w, r, result)
}

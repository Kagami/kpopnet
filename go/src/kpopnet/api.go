package kpopnet

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	maxOverheadSize = int64(10 * 1024)
	maxFileSize     = int64(5 * 1024 * 1024)
	maxBodySize     = maxFileSize + maxOverheadSize

	errInternal     = errors.New("internal error")
	errParseForm    = errors.New("error parsing form")
	errParseFile    = errors.New("error parsing form file")
	errRecognize    = errors.New("cannot recognize")
	errNoSingleFace = errors.New("not a single face")
)

func setApiHeaders(w http.ResponseWriter) {
	head := w.Header()
	head.Set("Cache-Control", "no-cache")
	head.Set("Content-Type", "application/json")
}

func serveData(w http.ResponseWriter, r *http.Request, data []byte) {
	etag := fmt.Sprintf("\"%s\"", hashBytes(data))
	if checkEtag(w, r, etag) {
		return
	}
	setApiHeaders(w)
	w.Header().Set("ETag", etag)
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
	// FIXME(Kagami): For some reason cached request is not fast enough.
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
	fd, err := fhs[0].Open()
	if err != nil {
		serveError(w, r, errParseFile, 400)
		return
	}
	defer fd.Close()
	fdata, err := ioutil.ReadAll(fd)
	if err != nil {
		serveError(w, r, errParseFile, 400)
		return
	}
	idolId, err := Recognize(fdata)
	if err != nil {
		serveError(w, r, errRecognize, 500)
		return
	}
	if idolId == nil {
		serveError(w, r, errNoSingleFace, 400)
		return
	}
	result := map[string]string{"id": *idolId}
	serveJson(w, r, result)
}

func Recognize(fdata []byte) (idolId *string, err error) {
	return
}

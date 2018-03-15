package kpopnet

import (
	"crypto/md5"
	"encoding/base64"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/dimfeld/httptreemux"
)

func logError(err error) {
	log.Printf("kpopnet: %s\n%s\n", err, debug.Stack())
}

func hashBytes(buf []byte) string {
	hash := md5.Sum(buf)
	return base64.RawStdEncoding.EncodeToString(hash[:])
}

func checkEtag(w http.ResponseWriter, r *http.Request, etag string) bool {
	if r.Header.Get("If-None-Match") == etag {
		w.WriteHeader(304)
		return true
	}
	return false
}

func getParam(r *http.Request, id string) string {
	return httptreemux.ContextParams(r.Context())[id]
}

package kpopnet

import (
	"net/http"
	"path/filepath"

	"github.com/dimfeld/httptreemux"
)

type ServerOptions struct {
	Address string
	WebRoot string
}

func StartServer(o ServerOptions) (err error) {
	router, err := createRouter(o)
	if err != nil {
		return
	}
	return http.ListenAndServe(o.Address, router)
}

func createRouter(o ServerOptions) (h http.Handler, err error) {
	r := httptreemux.NewContextMux()

	webRoot, err := filepath.Abs(o.WebRoot)
	if err != nil {
		return
	}
	indexPath := filepath.Join(webRoot, "index.html")
	faviconPath := filepath.Join(webRoot, "favicon.ico")
	staticRoot := filepath.Join(webRoot, "static")

	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(w, r, indexPath)
	})
	r.GET("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, faviconPath)
	})
	r.Handler("GET", "/static/*", http.StripPrefix("/static/",
		http.FileServer(http.Dir(staticRoot))))

	r.GET("/api/idols/profiles", ServeProfiles)
	r.POST("/api/idols/recognize", ServeRecognize)
	r.GET("/api/idols/by-image/:id", ServeImageInfo)

	h = http.Handler(r)
	return
}

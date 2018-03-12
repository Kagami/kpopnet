package server

import (
	"net/http"
	"path/filepath"

	"github.com/dimfeld/httptreemux"
)

type Options struct {
	Address string
	WebRoot string
}

func Start(o Options) (err error) {
	router, err := createRouter(o)
	if err != nil {
		return
	}
	return http.ListenAndServe(o.Address, router)
}

func createRouter(o Options) (h http.Handler, err error) {
	r := httptreemux.NewContextMux()

	webRoot, err := filepath.Abs(o.WebRoot)
	if err != nil {
		return
	}
	indexPath := filepath.Join(webRoot, "index.html")
	faviconPath := filepath.Join(webRoot, "favicon.ico")
	staticRoot := filepath.Join(webRoot, "static")

	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, indexPath)
	})
	r.GET("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, faviconPath)
	})
	r.Handler("GET", "/static/*", http.StripPrefix("/static/",
		http.FileServer(http.Dir(staticRoot))))

	h = http.Handler(r)
	return
}

package server

import (
	"net/http"
	"path/filepath"

	"github.com/dimfeld/httptreemux"
)

var (
	idolApi string
)

type Options struct {
	Address string
	WebRoot string
	IdolApi string
}

func Start(o Options) (err error) {
	idolApi = o.IdolApi
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
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(w, r, indexPath)
	})
	r.GET("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, faviconPath)
	})
	r.Handler("GET", "/static/*", http.StripPrefix("/static/",
		http.FileServer(http.Dir(staticRoot))))

	api := r.NewGroup("/api")
	api.GET("/profiles", serveProfiles)
	api.GET("/idols/*path", serveIdolApi)

	h = http.Handler(r)
	return
}

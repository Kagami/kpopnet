package server

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/dimfeld/httptreemux"
)

type Options struct {
	Address string
	WebRoot string
}

func Start(o Options) {
	router := createRouter(o)
	log.Printf("Listening on %v", o.Address)
	log.Fatal(http.ListenAndServe(o.Address, router))
}

func createRouter(o Options) http.Handler {
	r := httptreemux.NewContextMux()

	webRoot, _ := filepath.Abs(o.WebRoot)
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

	h := http.Handler(r)
	return h
}

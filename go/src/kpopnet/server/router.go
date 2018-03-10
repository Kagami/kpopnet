package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dimfeld/httptreemux"
)

const (
	DEFAULT_HOST = "127.0.0.1"
	DEFAULT_PORT = 8002
)

func Start() {
	addr := fmt.Sprintf("%v:%v", DEFAULT_HOST, DEFAULT_PORT)
	router := createRouter()
	log.Printf("Listening on %v", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func createRouter() http.Handler {
	r := httptreemux.NewContextMux()
	h := http.Handler(r)
	return h
}

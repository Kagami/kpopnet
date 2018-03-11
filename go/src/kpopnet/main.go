package main

import (
	"flag"

	"kpopnet/server"
)

func main() {
	var (
		address string
		webRoot string
	)
	flag.StringVar(
		&address,
		"b",
		"127.0.0.1:8002",
		"address to listen on for incoming HTTP connections",
	)
	flag.StringVar(
		&webRoot,
		"w",
		"./dist",
		"site dist directory location",
	)
	flag.Parse()

	opts := server.Options{
		Address: address,
		WebRoot: webRoot,
	}
	server.Start(opts)
}

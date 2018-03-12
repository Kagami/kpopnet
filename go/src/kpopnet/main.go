package main

import (
	"flag"
	"log"

	"kpopnet/db"
	"kpopnet/server"
)

func main() {
	var (
		address string
		webRoot string
		connStr string
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
		"site directory location",
	)
	// Designed to be integrated into meguca (cutechan) therefore the
	// default connection parameters choice.
	flag.StringVar(
		&connStr,
		"c",
		"user=meguca password=meguca dbname=meguca sslmode=disable",
		"PostgreSQL connection string",
	)
	flag.Parse()

	err := db.Start(connStr)
	if err != nil {
		log.Fatal(err)
	}

	serverOpts := server.Options{
		Address: address,
		WebRoot: webRoot,
	}
	log.Printf("Listening on %v", address)
	log.Fatal(server.Start(serverOpts))
}

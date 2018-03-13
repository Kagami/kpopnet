package main

import (
	"fmt"
	"log"

	"kpopnet/db"
	"kpopnet/profile"
	"kpopnet/server"

	"github.com/docopt/docopt-go"
)

const VERSION = "0.0.0"
const USAGE = `
K-pop neural network backend.

Usage:
  kpopnetd profile import [options]
  kpopnetd serve [options]
  kpopnetd [-h | --help]
  kpopnetd [-V | --version]

Options:
  -h --help     Show this screen.
  -V --version  Show version.
  -H <host>     Host to listen on [default: 127.0.0.1].
  -p <port>     Port to listen on [default: 8002].
  -c <conn>     PostgreSQL connection string
                [default: user=meguca password=meguca dbname=meguca sslmode=disable].
  -s <sitedir>  Site directory location [default: ./dist].
  -d <datadir>  Data directory location [default: ./data].
`

type config struct {
	Profile bool
	Import  bool
	Serve   bool
	Host    string `docopt:"-H"`
	Port    int    `docopt:"-p"`
	Conn    string `docopt:"-c"`
	Sitedir string `docopt:"-s"`
	Datadir string `docopt:"-d"`
}

func importProfiles(conf config) {
	err := db.Start(conf.Conn)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Importing profiles from %s", conf.Datadir)
	ps, err := profile.ReadAll(conf.Datadir)
	if err != nil {
		err = fmt.Errorf("Error reading profiles: %v", err)
		log.Fatal(err)
	}
	err = db.UpdateProfiles(ps)
	if err != nil {
		err = fmt.Errorf("Error updating DB profiles: %v", err)
		log.Fatal(err)
	}
	log.Print("Done.")
}

func serve(conf config) {
	err := db.Start(conf.Conn)
	if err != nil {
		log.Fatal(err)
	}
	opts := server.Options{
		Address: fmt.Sprintf("%v:%v", conf.Host, conf.Port),
		WebRoot: conf.Sitedir,
	}
	log.Printf("Listening on %v", opts.Address)
	log.Fatal(server.Start(opts))
}

func main() {
	opts, err := docopt.ParseArgs(USAGE, nil, VERSION)
	if err != nil {
		log.Fatal(err)
	}
	var conf config
	if err := opts.Bind(&conf); err != nil {
		log.Fatal(err)
	}

	if conf.Profile && conf.Import {
		importProfiles(conf)
	} else if conf.Serve {
		serve(conf)
	} else {
		log.Fatal("No command selected, try --help.")
	}
}

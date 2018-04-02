package main

import (
	"fmt"
	"log"

	"kpopnet"

	"github.com/docopt/docopt-go"
)

const VERSION = "0.0.0"
const USAGE = `
K-pop neural network backend.

Usage:
  kpopnetd profile import [options]
  kpopnetd image import [options]
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
	Image   bool
	Serve   bool
	Host    string `docopt:"-H"`
	Port    int    `docopt:"-p"`
	Conn    string `docopt:"-c"`
	SiteDir string `docopt:"-s"`
	DataDir string `docopt:"-d"`
}

func importProfiles(conf config) {
	log.Printf("Importing profiles from %s", conf.DataDir)
	if err := kpopnet.ImportProfiles(conf.Conn, conf.DataDir); err != nil {
		log.Fatal(err)
	}
	log.Print("Done.")
}

func importImages(conf config) {
	log.Printf("Importing images from %s", conf.DataDir)
	if err := kpopnet.ImportImages(conf.Conn, conf.DataDir); err != nil {
		log.Fatal(err)
	}
	log.Print("Done.")
}

func serve(conf config) {
	err := kpopnet.StartDb(nil, conf.Conn)
	if err != nil {
		log.Fatal(err)
	}
	opts := kpopnet.ServerOptions{
		Address: fmt.Sprintf("%v:%v", conf.Host, conf.Port),
		WebRoot: conf.SiteDir,
	}
	log.Printf("Listening on %v", opts.Address)
	log.Fatal(kpopnet.StartServer(opts))
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
	} else if conf.Image && conf.Import {
		importImages(conf)
	} else if conf.Serve {
		serve(conf)
	} else {
		log.Fatal("No command selected, try --help.")
	}
}

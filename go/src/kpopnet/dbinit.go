//go:generate go-bindata -o bin_data.go --pkg kpopnet --nometadata --prefix sql sql/...

package kpopnet

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

var (
	db       *sql.DB
	prepared = make(map[string]*sql.Stmt)
)

func getQuery(id string) string {
	name := id + ".sql"
	return string(MustAsset(name))
}

func prepare() (err error) {
	names := AssetNames()
	for _, name := range names {
		id := strings.TrimSuffix(name, ".sql")
		switch {
		case strings.HasPrefix(name, "init_"):
			// Do nothing.
		case strings.HasPrefix(name, "fn_"):
			if err = execQ(id); err != nil {
				return fmt.Errorf("error preparing %s: %v", name, err)
			}
		default:
			if prepared[id], err = db.Prepare(getQuery(id)); err != nil {
				return fmt.Errorf("error preparing %s: %v", name, err)
			}
		}
	}
	return
}

func StartDb(openedDb *sql.DB, connStr string) (err error) {
	if openedDb == nil {
		if db, err = sql.Open("postgres", connStr); err != nil {
			return
		}
	} else {
		db = openedDb
	}

	if err = execQ("init_db"); err != nil {
		return fmt.Errorf("error initializing database: %v", err)
	}

	if err = prepare(); err != nil {
		return
	}

	return
}

// Useful helper functions for working with database.
package db

import (
	"database/sql"
)

func exec(queryId string) (err error) {
	_, err = db.Exec(getQuery(queryId))
	return
}

func getTx() (tx *sql.Tx, err error) {
	return db.Begin()
}

func getRoTx() (tx *sql.Tx, err error) {
	tx, err = getTx()
	if err != nil {
		return
	}
	_, err = tx.Exec("SET TRANSACTION READ ONLY")
	if err != nil {
		err = tx.Rollback()
	}
	return
}

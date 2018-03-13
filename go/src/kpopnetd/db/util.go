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

func setReadOnly(tx *sql.Tx) (err error) {
	_, err = tx.Exec("SET TRANSACTION READ ONLY")
	return
}

func setRepeatableRead(tx *sql.Tx) (err error) {
	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")
	return
}

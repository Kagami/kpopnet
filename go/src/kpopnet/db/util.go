// Useful helper functions for working with database.
package db

import (
	"database/sql"
)

func execQ(queryId string) (err error) {
	_, err = db.Exec(getQuery(queryId))
	return
}

func beginTx() (tx *sql.Tx, err error) {
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

func endTx(tx *sql.Tx, err *error) {
	if *err != nil {
		tx.Rollback()
		return
	}
	*err = tx.Commit()
}

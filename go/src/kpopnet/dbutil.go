package kpopnet

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

func endTx(tx *sql.Tx, err *error) {
	if *err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			// Can only log this because original err should be preserved.
			logError(rbErr)
		}
		return
	}
	*err = tx.Commit()
}

func setReadOnly(tx *sql.Tx) (err error) {
	_, err = tx.Exec("SET TRANSACTION READ ONLY")
	return
}

func setRepeatableRead(tx *sql.Tx) (err error) {
	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")
	return
}

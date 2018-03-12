// Useful helper functions for working with database.
package db

func exec(queryId string) (err error) {
	_, err = db.Exec(getQuery(queryId))
	return
}

func scanBool(queryId string, args ...interface{}) (val bool, err error) {
	err = db.QueryRow(getQuery(queryId), args...).Scan(&val)
	return
}

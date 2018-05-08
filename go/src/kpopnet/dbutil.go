package kpopnet

import (
	"database/sql"
	"fmt"
	"image"
	"regexp"
	"strconv"
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

// PostgreSQL to Go type mappers.

func rect2str(rect image.Rectangle) string {
	x0 := rect.Min.X
	y0 := rect.Min.Y
	x1 := rect.Max.X
	y1 := rect.Max.Y
	return fmt.Sprintf("((%d,%d),(%d,%d))", x0, y0, x1, y1)
}

var rectRe = regexp.MustCompile(`^\((\d+),(\d+)\),\((\d+),(\d+)\)$`)

func str2rect(rectStr string) (rect image.Rectangle) {
	coords := rectRe.FindStringSubmatch(rectStr)
	// Shouldn't happen because PostgreSQL rectangle format is fixed so
	// return just an empty rect in case of code mistake.
	if coords == nil {
		return
	}
	x0, _ := strconv.Atoi(coords[1])
	y0, _ := strconv.Atoi(coords[2])
	x1, _ := strconv.Atoi(coords[3])
	y1, _ := strconv.Atoi(coords[4])
	return image.Rect(x0, y0, x1, y1)
}

package tinydb

import "database/sql"

func Dui(db *sql.DB, sql string, args ...interface{}) (r sql.Result, err error) {
	tx, err := db.Begin()
	if err != nil {
		return r, err
	}
	if len(args) > 0 {
		r, err = tx.Exec(sql, args...)
	} else {
		r, err = tx.Exec(sql)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return r, err
	}
	return
}

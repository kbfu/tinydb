package tinydb

import "database/sql"

func Dui(db *sql.DB, sql string, args ...string) (r sql.Result, err error) {
	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		return r, err
	}
	r, err = tx.Exec(sql, args)
	if err != nil {
		return r, err
	}
	err = tx.Commit()
	if err != nil {
		return r, err
	}
	return
}

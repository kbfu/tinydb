package tinydb

import (
	"database/sql"
	"fmt"
)

type Exec struct {
	err  error
	db   TinyDb
	sql  string
	args []string
}

func (db *TinyDb) Exec(sql string, args ...interface{}) (result sql.Result, err error) {
	if db.Debug {
		fmt.Println(db.sqlDb)
		fmt.Println(sql, args)
	}
	result, err = Dui(db.sqlDb, sql, args...)
	return
}

package tinydb

import (
	"database/sql"
	"fmt"
	"reflect"
)

type M map[string]interface{}

func Where(condition M) string {
	var where []string
	for k, v := range condition {
		switch reflect.ValueOf(v).Kind() {
		case reflect.Int:
			where = append(where, fmt.Sprintf("%s = %v", k, v))
		case reflect.String:
			where = append(where, fmt.Sprintf("%s = \"%v\"", k, v))
		}

	}
	whereStr := ""
	for k, v := range where {
		if k != len(where)-1 {
			whereStr = whereStr + v + " AND "
		} else {
			whereStr = whereStr + v
		}
	}
	return "WHERE " + whereStr
}

func From(table interface{}) (from string) {
	switch reflect.TypeOf(table).Kind() {
	case reflect.String:
		from = "FROM " + table.(string)
	case reflect.Ptr:
		from = "FROM (" + fmt.Sprintf("SELECT %s %s %s",
			table.(*Select).columns, table.(*Select).from, table.(*Select).where) + ")"
	}
	return
}

func Dui(db *sql.DB, sql string) (r sql.Result, err error) {
	tx, err := db.Begin()
	defer tx.Rollback()
	if err != nil {
		return r, err
	}
	r, err = tx.Exec(sql)
	if err != nil {
		return r, err
	}
	err = tx.Commit()
	if err != nil {
		return r, err
	}
	return
}

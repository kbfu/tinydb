package tinydb

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Query struct {
	rows *sql.Rows
	err  error
	db   TinyDb
	sql  string
	args []interface{}
}

func (db *TinyDb) Query(sql string, args ...interface{}) *Query {
	var q Query
	q.db = *db
	q.sql = sql
	q.args = args
	if len(args) > 0 {
		q.rows, q.err = db.sqlDb.Query(sql, args)
	} else {
		q.rows, q.err = db.sqlDb.Query(sql)
	}
	return &q
}

func (q *Query) Get(obj interface{}) (err error) {
	if q.err != nil {
		return q.err
	}
	if q.db.Debug {
		fmt.Println(q.db.sqlDb)
		fmt.Println(fmt.Sprintf(q.sql, q.args))
	}
	defer q.rows.Close()
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return errors.New("not pointer type")
	}
	v = v.Elem()
	for q.rows.Next() {
		var fields []interface{}
		var elem reflect.Value
		columns, _ := q.rows.Columns()
		switch v.Kind() {
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				tag := v.Type().Field(i).Tag.Get("db")
				if tag == "" {
					tag = v.Type().Field(i).Name
				}
				for _, c := range columns {
					if strings.ToLower(tag) == strings.ToLower(c) {
						fields = append(fields, v.Field(i).Addr().Interface())
					}
				}
			}
		case reflect.Slice:
			elem = reflect.New(v.Type().Elem()).Elem()
			if elem.Kind() == reflect.Struct {
				for i := 0; i < elem.NumField(); i++ {
					tag := elem.Type().Field(i).Tag.Get("db")
					if tag == "" {
						tag = v.Type().Field(i).Name
					}
					for _, c := range columns {
						if strings.ToLower(tag) == strings.ToLower(c) {
							fields = append(fields, elem.Field(i).Addr().Interface())
						}
					}
				}
			} else {
				fields = make([]interface{}, 1)
				fields[0] = elem.Addr().Interface()
			}
		default:
			fields = make([]interface{}, 1)
			fields[0] = v.Addr().Interface()
		}
		// check tag
		fmt.Println(q.rows.Columns())
		err := q.rows.Scan(fields...)
		if err != nil {
			return err
		}
		if v.Kind() == reflect.Slice {
			v.Set(reflect.Append(v, elem))
		}
	}
	err = q.rows.Err()
	if err != nil {
		return
	}
	return
}

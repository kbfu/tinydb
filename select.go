package tinydb

import (
	"errors"
	"fmt"
	"reflect"
)

type Select struct {
	db      *TinyDb
	columns string
	from    interface{}
	where   string
	as      string
}

func (db *TinyDb) Select(columns ...string) *Select {
	var s Select
	s.db = db
	cols := ""
	for k, v := range columns {
		if k != len(columns)-1 {
			cols = cols + v + ","
		} else {
			cols = cols + v
		}
	}
	s.columns = cols
	return &s
}

func (s *Select) From(table interface{}) *Select {
	s.from = From(table)
	return s
}

func (s *Select) Where(condition M) *Select {
	s.where = Where(condition)
	return s
}

func (s *Select) As(alias string) *Select {
	s.as = "AS " + alias
	return s
}

func (s *Select) Exec(obj interface{}) (err error) {
	rows, err := s.db.Query(fmt.Sprintf("SELECT %s %s %s %s", s.columns, s.from, s.where, s.as))
	if err != nil {
		return err
	}
	defer rows.Close()
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return errors.New("not pointer type")
	}
	v = v.Elem()
	for rows.Next() {
		var fields []interface{}
		var elem reflect.Value
		switch v.Kind() {
		case reflect.Struct:
			fields = make([]interface{}, v.NumField())
			for i := 0; i < v.NumField(); i++ {
				fields[i] = v.Field(i).Addr().Interface()
			}
		case reflect.Slice:
			elem = reflect.New(v.Type().Elem()).Elem()
			if elem.Kind() == reflect.Struct {
				fields = make([]interface{}, elem.NumField())
				for i := 0; i < elem.NumField(); i++ {
					fields[i] = elem.Field(i).Addr().Interface()
				}
			} else {
				fields = make([]interface{}, 1)
				fields[0] = elem.Addr().Interface()
			}
		default:
			fields = make([]interface{}, 1)
			fields[0] = v.Addr().Interface()
		}
		err := rows.Scan(fields...)
		if err != nil {
			return err
		}
		if v.Kind() == reflect.Slice {
			v.Set(reflect.Append(v, elem))
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return
}

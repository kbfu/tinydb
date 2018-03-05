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
	err     error
}

type Columner interface {
	Cols() string
}

type Counter interface {
	Count() string
	Columner
}

type CountType struct {
	Col string
	ColumnType
}

type ColumnType struct {
	ColNames []string
}

func (c *CountType) Count() (col string) {
	if c.Col == "*" {
		col = fmt.Sprintf("count(%s)", c.Col)
	} else {
		col = fmt.Sprintf("count(`%s`)", c.Col)
	}
	return
}

func Count(col string) *CountType {
	c := CountType{}
	c.Col = col
	return &c
}

func (c *ColumnType) Cols() string {
	cols := ""
	for k, v := range c.ColNames {
		if k != len(c.ColNames)-1 {
			cols = cols + fmt.Sprintf("`%s`,", v)
		} else {
			cols = cols + fmt.Sprintf("`%s`", v)
		}
	}
	return cols
}

func Columns(cols ...string) *ColumnType {
	c := ColumnType{}
	for _, v := range cols {
		c.ColNames = append(c.ColNames, v)
	}
	return &c
}

func (db *TinyDb) Select(columner Columner) *Select {
	var (
		s    Select
		cols string
	)
	s.db = db
	switch columner.(type) {
	case Counter:
		cols = columner.(Counter).Count()
	case Columner:
		cols = columner.Cols()
	default:
		s.err = errors.New("not implemented")
	}
	s.columns = cols
	return &s
}

func (s *Select) From(table interface{}) *Select {
	s.from = From(table)
	return s
}

func (s *Select) Where(condition ...WhereConditioner) *Select {
	s.where, s.err = Where(condition...)
	return s
}

func (s *Select) As(alias string) *Select {
	s.as = "AS " + alias
	return s
}

func (s *Select) Exec(obj interface{}) (err error) {
	if s.err != nil {
		return s.err
	}
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
		return
	}
	return
}

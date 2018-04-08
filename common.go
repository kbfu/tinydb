package tinydb

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

type M map[string]interface{}

type EqualStruct struct {
	M
}

type LikeStruct struct {
	M
}

type InStruct struct {
	M
}

func (i *InStruct) Where() (wheres []string, err error) {
	for k, v := range i.M {
		in := fmt.Sprintf("%s in (", k)
		val := reflect.ValueOf(v)
		switch val.Kind() {
		case reflect.Slice:
			switch reflect.New(val.Type().Elem()).Elem().Kind() {
			case reflect.String:
				for k, val := range v.([]string) {
					if k == len(v.([]string))-1 {
						in += fmt.Sprintf("'%s')", val)
					} else {
						in += fmt.Sprintf("'%s',", val)
					}
				}
				wheres = append(wheres, in)
			case reflect.Int:
				for k, val := range v.([]int) {
					if k == len(v.([]int))-1 {
						in += fmt.Sprintf("'%v')", val)
					} else {
						in += fmt.Sprintf("'%v',", val)
					}
				}
				wheres = append(wheres, in)
			case reflect.Int64:
				for k, val := range v.([]int64) {
					if k == len(v.([]int64))-1 {
						in += fmt.Sprintf("'%v')", val)
					} else {
						in += fmt.Sprintf("'%v',", val)
					}
				}
				wheres = append(wheres, in)
			default:
				return wheres, errors.New("type not supported")
			}
		default:
			return wheres, errors.New("type not supported")
		}
	}
	return
}

func (e *EqualStruct) Where() (wheres []string, err error) {
	for k, v := range e.M {
		switch reflect.ValueOf(v).Kind() {
		case reflect.Int, reflect.Int64, reflect.String:
			wheres = append(wheres, fmt.Sprintf("`%s` = '%v'", k, v))
		default:
			return wheres, errors.New("type not supported")
		}
	}
	return
}

func Like(pairs M) *LikeStruct {
	return &LikeStruct{pairs}
}

func Equal(pairs M) *EqualStruct {
	return &EqualStruct{pairs}
}

func In(pairs M) *InStruct {
	return &InStruct{pairs}
}

func (l *LikeStruct) Where() (wheres []string, err error) {
	for k, v := range l.M {
		switch reflect.ValueOf(v).Kind() {
		case reflect.String, reflect.Int:
			wheres = append(wheres, fmt.Sprintf("`%s` like '%v'", k, v))
		default:
			return wheres, errors.New("type not supported")
		}
	}
	return
}

type WhereConditioner interface {
	Where() (wheres []string, err error)
}

func Where(condition ...WhereConditioner) (where string, err error) {
	var wheres []string
	for _, v := range condition {
		pairs, err := v.Where()
		if err != nil {
			return where, err
		}
		wheres = append(wheres, pairs...)
	}
	whereStr := ""
	for k, v := range wheres {
		if k != len(wheres)-1 {
			whereStr = whereStr + v + " AND "
		} else {
			whereStr = whereStr + v
		}
	}
	return "WHERE " + whereStr, nil
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

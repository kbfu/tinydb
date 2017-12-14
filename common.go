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

func (e *EqualStruct) Where() (wheres []string, err error) {
	for k, v := range e.M {
		switch reflect.ValueOf(v).Kind() {
		case reflect.Int:
			wheres = append(wheres, fmt.Sprintf("%s = %v", k, v))
		case reflect.String:
			wheres = append(wheres, fmt.Sprintf("%s = \"%v\"", k, v))
		default:
			return wheres, errors.New("type not supported")
		}
	}
	return
}

func Like(pairs M) *LikeStruct {
	return &LikeStruct{pairs}
}

func Equal(paris M) *EqualStruct {
	return &EqualStruct{paris}
}

func (l *LikeStruct) Where() (wheres []string, err error) {
	for k, v := range l.M {
		switch reflect.ValueOf(v).Kind() {
		case reflect.String, reflect.Int:
			wheres = append(wheres, fmt.Sprintf("%s like \"%v\"", k, v))
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

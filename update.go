package tinydb

import (
	"fmt"
)

type Update struct {
	db    *TinyDb
	set   string
	table string
	where string
	err   error
}

func (db *TinyDb) Update(table string) *Update {
	var u Update
	u.db = db
	u.table = table
	return &u
}

func (u *Update) Where(condition ...WhereConditioner) *Update {
	u.where, u.err = Where(condition...)
	return u
}

func (u *Update) Set(condition M) *Update {
	var set []string
	for k, v := range condition {
		set = append(set, fmt.Sprintf("%s = \"%s\"", k, v))
	}
	setStr := ""
	for k, v := range set {
		if k != len(set)-1 {
			setStr = setStr + v + ", "
		} else {
			setStr = setStr + v
		}
	}
	u.set = "SET " + setStr
	return u
}

func (u *Update) Exec() (err error) {
	_, err = Dui(u.db.DB, fmt.Sprintf("UPDATE %s %s %s", u.table, u.set, u.where))
	if err != nil {
		return err
	}
	return
}

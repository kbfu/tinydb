package tinydb

import (
	"fmt"
	"strings"
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
		switch v.(type) {
		case string:
			val := strings.Replace(v.(string), "\\n", "\\\\n",-1)
			val = strings.Replace(val, "\"", "\\\"", -1)
			val = strings.Replace(val, "'", "\\'", -1)
			set = append(set, fmt.Sprintf("`%s` = '%s'", k, val))
		default:
			set = append(set, fmt.Sprintf("`%s` = '%v'", k, v))
		}

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
	if u.db.Debug {
		fmt.Println(u.db.sqlDb)
		fmt.Println(fmt.Sprintf("UPDATE %s %s %s", u.table, u.set, u.where))
	}
	_, err = Dui(u.db.sqlDb, fmt.Sprintf("UPDATE %s %s %s", u.table, u.set, u.where))
	if err != nil {
		return err
	}
	if u.err != nil {
		return u.err
	}
	return
}

package tinydb

import (
	"fmt"
	"strings"
)

type Insert struct {
	db      *TinyDb
	table   string
	columns string
	values  string
}

func (db *TinyDb) Insert() *Insert {
	var i Insert
	i.db = db
	return &i
}

func (i *Insert) Into(table string) *Insert {
	i.table = table
	return i
}

func (i *Insert) Columns(columns ...string) *Insert {
	cols := ""
	if len(columns) > 0 {
		cols = "("
	}
	for k, v := range columns {
		if k != len(columns)-1 {
			cols = cols + fmt.Sprintf("`%s`,", v)
		} else {
			cols = cols + fmt.Sprintf("`%s`)", v)
		}
	}
	i.columns = cols
	return i
}

func (i *Insert) Values(values ...interface{}) *Insert {
	vals := "VALUES ("
	for k, v := range values {
		if k != len(values)-1 {
			switch v.(type) {
			case string:
				vals = vals + fmt.Sprintf("'%s'", strings.Replace(v.(string), "\\n", "\\\\n", -1)) + ","
				vals = strings.Replace(vals, "\\\"", "\\\\\"", -1)
			default:
				vals = vals + fmt.Sprintf("'%v'", v) + ","
			}
		} else {
			switch v.(type) {
			case string:
				vals = vals + fmt.Sprintf("'%s'", strings.Replace(v.(string), "\\n", "\\\\n", -1)) + ")"
				vals = strings.Replace(vals, "\\\"", "\\\\\"", -1)
			default:
				vals = vals + fmt.Sprintf("'%v'", v) + ")"
			}
		}
	}
	i.values = vals
	return i
}

func (i *Insert) Exec() (id int64, err error) {
	r, err := Dui(i.db.sqlDb, fmt.Sprintf("INSERT INTO %s %s %s", i.table, i.columns, i.values))
	if err != nil {
		return id, err
	}
	id, err = r.LastInsertId()
	if err != nil {
		return id, err
	}
	return
}

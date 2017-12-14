package tinydb

import "fmt"

type Delete struct {
	db    *TinyDb
	from  string
	where string
}

func (db *TinyDb) Delete() *Delete {
	var d Delete
	d.db = db
	return &d
}

func (d *Delete) From(table string) *Delete {
	d.from = From(table)
	return d
}

func (d *Delete) Where(condition M) *Delete {
	d.where = Where(condition)
	return d
}

func (d *Delete) Exec() (err error) {
	_, err = Dui(d.db.DB, fmt.Sprintf("DELETE %s %s", d.from, d.where))
	if err != nil {
		return err
	}
	return
}

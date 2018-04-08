package tinydb

import (
	"fmt"
	"log"
)

type Delete struct {
	db    *TinyDb
	from  string
	where string
	err   error
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

func (d *Delete) Where(condition ...WhereConditioner) *Delete {
	d.where, d.err = Where(condition...)
	return d
}

func (d *Delete) Exec() (err error) {
	if d.db.Debug {
		log.Println(d.db.sqlDb)
		log.Println(fmt.Sprintf("DELETE %s %s", d.from, d.where))
	}
	_, err = Dui(d.db.sqlDb, fmt.Sprintf("DELETE %s %s", d.from, d.where))
	if err != nil {
		return err
	}
	if d.err != nil {
		return d.err
	}
	return
}

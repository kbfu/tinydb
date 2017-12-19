package tinydb

import (
	"database/sql"
	"fmt"
)

type TinyDb struct {
	*sql.DB
}

func New(dbType, user, password, host, port, database, charset string) (db TinyDb, err error) {
	sqlDb, err := sql.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", user, password, host, port,
		database, charset))
	if err != nil {
		return db, err
	}
	db.DB = sqlDb
	return db, err
}

func (db *TinyDb) SetMaxIdleConns(n int) {
	db.DB.SetMaxIdleConns(n)
}

func (db *TinyDb) SetMaxOpenConns(n int) {
	db.DB.SetMaxOpenConns(n)
}

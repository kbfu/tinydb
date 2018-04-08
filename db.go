package tinydb

import (
	"database/sql"
	"fmt"
	"log"
)

type TinyDb struct {
	sqlDb *sql.DB
	Debug bool
}

func New(dbType, user, password, host, port, database, charset string) (db TinyDb, err error) {
	sqlDb, err := sql.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", user, password, host, port,
		database, charset))
	if err != nil {
		return db, err
	}
	db.sqlDb = sqlDb
	return db, err
}

func (db *TinyDb) SetMaxIdleConns(n int) {
	db.sqlDb.SetMaxIdleConns(n)
}

func (db *TinyDb) SetMaxOpenConns(n int) {
	db.sqlDb.SetMaxOpenConns(n)
}

func (db *TinyDb) SetDebug() {
	db.Debug = true
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func (db *TinyDb) Ping() error {
	err := db.sqlDb.Ping()
	return err
}

package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

func (db *DB) AddPath(path string) error {
	stmt, err := db.conn.Prepare("INSERT INTO pathinfo(prio, path, timestamp) values(?, ?, datetime())")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(1, path)
	if err != nil {
		return err
	}

	return nil
}

func initDB(db *sql.DB) error {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS pathinfo(id INTEGER PRIMARY KEY, prio INTEGER, path TEXT, timestamp DATETIME)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func openDB(path string) (*DB, error) {
	ret := new(DB)
	var err error

	ret.conn, err = sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	err = initDB(ret.conn)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

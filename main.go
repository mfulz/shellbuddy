package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type args struct {
	config string
	add    bool
	path   bool
	cmd    bool
}

func parseArgs() (*args, error) {
	ret := new(args)

	flag.StringVar(&ret.config, "config", "~/.shellbuddy/config", "Configuration file to use. Defaults to ~/.shellbuddy/config")
	flag.BoolVar(&ret.add, "add", false, "If you want to add / update an entry")
	flag.BoolVar(&ret.cmd, "cmd", false, "If you want to add / update commands")
	flag.BoolVar(&ret.path, "path", false, "If you want to add / update pathes")

	flag.Parse()

	if !ret.cmd && !ret.path {
		return nil, fmt.Errorf("Neither path nor cmd selected")
	}

	return ret, nil
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

func main() {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		panic(err)
	}

	err = initDB(db)
	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare("INSERT INTO pathinfo(prio, path, timestamp) values(?, ?, datetime())")
	if err != nil {
		panic(err)
	}

	res, err := stmt.Exec(1, "/home/mfulz/Projects")
	if err != nil {
		panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Last id: %v\n", id)
}

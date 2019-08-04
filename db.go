package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Path struct {
	id       int64
	prio     int64
	path     string
	datetime *time.Time
}

type DB struct {
	conn             *sql.DB
	addPathStmt      *sql.Stmt
	selectPathStmt   *sql.Stmt
	updatePathStmt   *sql.Stmt
	selectPathesStmt *sql.Stmt
}

func (db *DB) AddPath(path string) error {
	e, err := db.GetPath(path)
	if err != nil {
		return nil
	}

	if e.id == -1 {
		_, err := db.addPathStmt.Exec(1, path)
		if err != nil {
			return err
		}
	} else {
		_, err := db.updatePathStmt.Exec(e.prio+1, e.id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) GetPath(path string) (Path, error) {
	ret := Path{-1, 0, "", nil}

	res, err := db.selectPathStmt.Query(path)
	if err != nil {
		return ret, err
	}
	defer res.Close()

	if !res.Next() {
		return ret, nil
	}

	err = res.Scan(&ret.id, &ret.prio, &ret.path, &ret.datetime)

	return ret, err
}

func (db *DB) GetPathes() ([]Path, error) {
	ret := []Path{}

	res, err := db.selectPathesStmt.Query()
	if err != nil {
		return ret, err
	}
	defer res.Close()

	for res.Next() {
		e := Path{}
		err := res.Scan(&e.id, &e.prio, &e.path, &e.datetime)
		if err != nil {
			return []Path{}, err
		}

		ret = append(ret, e)
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

	ret.addPathStmt, err = ret.conn.Prepare("INSERT INTO pathinfo(prio, path, timestamp) values(?, ?, datetime())")
	if err != nil {
		return nil, err
	}

	ret.selectPathStmt, err = ret.conn.Prepare("SELECT id, prio, path, timestamp FROM pathinfo where path =?")
	if err != nil {
		return nil, err
	}

	ret.selectPathesStmt, err = ret.conn.Prepare("SELECT id, prio, path, timestamp FROM pathinfo ORDER BY prio DESC, timestamp DESC")
	if err != nil {
		return nil, err
	}

	ret.updatePathStmt, err = ret.conn.Prepare("UPDATE pathinfo SET prio=?, timestamp=datetime() where id=?")
	if err != nil {
		return nil, err
	}

	return ret, nil
}

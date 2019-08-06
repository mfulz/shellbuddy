package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn                          *sql.DB
	addEntryStmt                  *sql.Stmt
	selectEntryStmt               *sql.Stmt
	updateEntryStmt               *sql.Stmt
	selectEntriesStmt             *sql.Stmt
	selectAllEntriesStmt          *sql.Stmt
	selectEntriesByStringStmt     *sql.Stmt
	selectAllEntriesByStringStmt  *sql.Stmt
	selectEntriesByTypeStmt       *sql.Stmt
	selectEntriesByStringTypeStmt *sql.Stmt
}

func (db *DB) getEntry(text string, etype EntryType) (Entry, error) {
	ret := Entry{-1, 0, "", nil, -1}

	res, err := db.selectEntryStmt.Query(text, etype)
	if err != nil {
		return ret, err
	}
	defer res.Close()

	if !res.Next() {
		return ret, nil
	}

	err = res.Scan(&ret.id, &ret.prio, &ret.text, &ret.timestamp, &ret.etype)

	return ret, err
}

func (db *DB) addEntry(text string, etype EntryType) error {
	e, err := db.getEntry(text, etype)
	if err != nil {
		return err
	}

	if e.id == -1 {
		_, err := db.addEntryStmt.Exec(1, text, etype)
		if err != nil {
			return err
		}
	} else {
		_, err := db.updateEntryStmt.Exec(e.prio+1, e.id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) getEntries(search string, etype EntryType) ([]Entry, error) {
	ret := []Entry{}

	var err error
	var res *sql.Rows

	if search == "" {
		res, err = db.selectEntriesStmt.Query(etype)
	} else {
		res, err = db.selectEntriesByStringStmt.Query("%"+search+"%", etype)
	}
	if err != nil {
		return ret, err
	}
	defer res.Close()

	for res.Next() {
		e := Entry{}
		err := res.Scan(&e.id, &e.prio, &e.text, &e.timestamp, &e.etype)
		if err != nil {
			return []Entry{}, err
		}

		ret = append(ret, e)
	}

	return ret, nil
}

func (db *DB) GetAllEntries(search string) ([]Entry, error) {
	ret := []Entry{}

	var err error
	var res *sql.Rows

	if search == "" {
		res, err = db.selectAllEntriesStmt.Query()
	} else {
		res, err = db.selectAllEntriesByStringStmt.Query("%" + search + "%")
	}
	if err != nil {
		return ret, err
	}
	defer res.Close()

	for res.Next() {
		e := Entry{}
		err := res.Scan(&e.id, &e.prio, &e.text, &e.timestamp, &e.etype)
		if err != nil {
			return []Entry{}, err
		}

		ret = append(ret, e)
	}

	return ret, nil
}

func (db *DB) GetAllEntriesByType(search string, etypes []EntryType) ([]Entry, error) {
	ret := []Entry{}

	var res *sql.Rows

	if search == "" {
		qetypes := make([]interface{}, len(etypes))
		for i, e := range etypes {
			qetypes[i] = e
		}

		stmt, err := db.genInQuery(len(qetypes), "SELECT id, prio, text, timestamp, etype FROM entries WHERE etype IN (", ") ORDER BY prio DESC, timestamp DESC")
		if err != nil {
			return nil, err
		}
		defer stmt.Close()

		res, err = stmt.Query(qetypes...)
		if err != nil {
			return ret, err
		}
		defer res.Close()
	} else {
		qetypes := make([]interface{}, len(etypes)+1)
		qetypes[0] = "%" + search + "%"
		for i, e := range etypes {
			qetypes[i+1] = e
		}

		stmt, err := db.genInQuery(len(qetypes), "SELECT id, prio, text, timestamp, etype FROM entries WHERE text LIKE ? AND etype IN (", ") ORDER BY prio DESC, timestamp DESC")
		if err != nil {
			return nil, err
		}
		defer stmt.Close()

		res, err = stmt.Query(qetypes...)
		if err != nil {
			return ret, err
		}
		defer res.Close()
	}

	for res.Next() {
		e := Entry{}
		err := res.Scan(&e.id, &e.prio, &e.text, &e.timestamp, &e.etype)
		if err != nil {
			return []Entry{}, err
		}

		ret = append(ret, e)
	}

	return ret, nil
}

func (db *DB) AddPath(path string) error {
	err := db.addEntry(path, DIR)
	return err
}

func (db *DB) GetPath(path string) (Entry, error) {
	ret, err := db.getEntry(path, DIR)
	return ret, err
}

func (db *DB) GetPathes(search string) ([]Entry, error) {
	ret, err := db.getEntries(search, DIR)
	return ret, err
}

func (db *DB) AddCmd(cmd string) error {
	err := db.addEntry(cmd, COMMAND)
	return err
}

func (db *DB) GetCmd(cmd string) (Entry, error) {
	ret, err := db.getEntry(cmd, COMMAND)
	return ret, err
}

func (db *DB) GetCmds(cmd string) ([]Entry, error) {
	ret, err := db.getEntries(cmd, COMMAND)
	return ret, err
}

func (db *DB) genInQuery(l int, queryStart string, queryEnd string) (*sql.Stmt, error) {
	ret, err := db.conn.Prepare(queryStart + genQueryPlaceholders(l) + queryEnd)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (db *DB) Close() {
	db.addEntryStmt.Close()
	db.selectEntryStmt.Close()
	db.updateEntryStmt.Close()
	db.selectEntriesStmt.Close()
	db.selectAllEntriesStmt.Close()
	db.selectEntriesByStringStmt.Close()
	db.selectAllEntriesByStringStmt.Close()
	db.conn.Close()
}

func initDB(db *sql.DB) error {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS entries(id INTEGER PRIMARY KEY, prio INTEGER, text TEXT, timestamp DATETIME, etype INTEGER)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func genQueryPlaceholders(l int) string {
	ret := ""

	for i := 0; i < l; i++ {
		if ret != "" {
			ret += ",?"
		} else {
			ret = "?"
		}
	}

	return ret
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

	ret.addEntryStmt, err = ret.conn.Prepare("INSERT INTO entries(prio, text, timestamp, etype) values(?, ?, datetime(), ?)")
	if err != nil {
		return nil, err
	}

	ret.selectEntryStmt, err = ret.conn.Prepare("SELECT id, prio, text, timestamp, etype FROM entries where text =? AND etype =?")
	if err != nil {
		return nil, err
	}

	ret.updateEntryStmt, err = ret.conn.Prepare("UPDATE entries SET prio=?, timestamp=datetime() where id=?")
	if err != nil {
		return nil, err
	}

	ret.selectEntriesStmt, err = ret.conn.Prepare("SELECT id, prio, text, timestamp, etype FROM entries WHERE etype =? ORDER BY prio DESC, timestamp DESC")
	if err != nil {
		return nil, err
	}

	ret.selectAllEntriesStmt, err = ret.conn.Prepare("SELECT id, prio, text, timestamp, etype FROM entries ORDER BY prio DESC, timestamp DESC")
	if err != nil {
		return nil, err
	}

	ret.selectEntriesByStringStmt, err = ret.conn.Prepare("SELECT id, prio, text, timestamp, etype FROM entries WHERE text LIKE ? AND etype =? ORDER BY prio DESC, timestamp DESC")
	if err != nil {
		return nil, err
	}

	ret.selectAllEntriesByStringStmt, err = ret.conn.Prepare("SELECT id, prio, text, timestamp, etype FROM entries WHERE text LIKE ? ORDER BY prio DESC, timestamp DESC")
	if err != nil {
		return nil, err
	}

	return ret, nil
}

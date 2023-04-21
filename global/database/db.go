package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDataBase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./backend.db")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return db, nil
}

func InitDB() bool {
	return CreateTable()
}

func CreateTable() bool {
	db, err := OpenDataBase()
	if err != nil {
		log.Println(err)
		return false
	}
	defer db.Close()

	addressTableSql := `
	CREATE TABLE IF NOT EXISTS address (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		address TEXT UNIQUE
	);
	`

	storageTableSql := `
	CREATE TABLE IF NOT EXISTS storage (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		address_id INTEGER NOT NULL,
		stype     INTEGER,
		pay_space INTEGER,
		free_space INTEGER,
		used_space INTEGER,
		file_count INTEGER,
		update_time DATETIME,
		UNIQUE (address_id, stype) ON CONFLICT IGNORE,
		FOREIGN KEY (address_id) REFERENCES address(id)
	);
	`

	pkgInfoTableSql := `
	CREATE TABLE IF NOT EXISTS pkginfo (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		address_id INTEGER NOT NULL,
		stype INTEGER,
		hash TEXT,
		size INTEGER,
		is_updated BOOLEAN,
		time DATETIME,
		FOREIGN KEY (address_id) REFERENCES address(id)
	);
	`

	ipfsObjectSql := `CREATE TABLE IF NOT EXISTS ipfsobjects (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        address_id INTEGER,
        name TEXT,
        size INTEGER,
        hash TEXT,
        FOREIGN KEY(address_id) REFERENCES address(id)
    );`

	_, err = db.Exec(addressTableSql)
	if err != nil {
		log.Println(err)
		return false
	}

	_, err = db.Exec(storageTableSql)
	if err != nil {
		log.Println(err)
		return false
	}

	_, err = db.Exec(pkgInfoTableSql)
	if err != nil {
		log.Println(err)
		return false
	}

	_, err = db.Exec(ipfsObjectSql)
	if err != nil {
		log.Println(err)
		return false
	}

	log.Println("Create Table Success!")
	return true
}

func execSql(db *sql.DB, message string, args ...any) error {
	stmt, err := db.Prepare(message)
	if err != nil {
		log.Println("exec sql error:", err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		log.Println("get result error: ", err)
		log.Println(args...)
		return err
	}
	return err
}

package global

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DataBase struct {
	DB *sql.DB
}

func InitDB() DataBase {
	db, err := sql.Open("sqlite3", "./backend.db")
	if err != nil {
		log.Println(err)
		return DataBase{}
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Println(err)
		return DataBase{}
	}
	DB := DataBase{db}
	log.Println("connect database success")
	DB.CreateTable()
	return DB
}

func (d *DataBase) CreateTable() {
	addressTableSql := `
	CREATE TABLE IF NOT EXISTS address (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		address TEXT
	);
	`

	_, err := d.DB.Exec(addressTableSql)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("create table success!")
}

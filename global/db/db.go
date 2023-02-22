package global

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB() {
	db, err := sql.Open("postgres", "postgres://gateway:Memo1234@127.0.0.1:5432/gateway?sslmode=disable")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("connect database success")
	DB = db
}
